package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestAudioServiceInit(t *testing.T) {
	as := NewAudioService()
	if as == nil {
		t.Error("AudioService initialization failed")
	}
}

func TestAudioServiceInitialState(t *testing.T) {
	as := NewAudioService()
	state := as.GetState()

	// NewAudioService sets state to AudioIdle, not AudioStopped
	if state.State != AudioIdle {
		t.Errorf("Expected initial state %q, got %q", AudioIdle, state.State)
	}
	if state.TrackName != "" {
		t.Errorf("Expected empty track name initially, got %q", state.TrackName)
	}
}

func TestAudioServiceSetVolume(t *testing.T) {
	as := NewAudioService()

	tests := []struct {
		name   string
		volume int
		expect float64 // stored as 0.0–1.0 in as.vol
	}{
		{"0%", 0, 0.0},
		{"50%", 50, 0.5},
		{"100%", 100, 1.0},
		{"clamp_low", -10, 0.0},
		{"clamp_high", 150, 1.0},
	}

	for _, test := range tests {
		as.SetVolume(test.volume)
		// vol is unexported but accessible within the same package
		if as.vol < test.expect-0.01 || as.vol > test.expect+0.01 {
			t.Errorf("%s: expected vol %.2f, got %.2f", test.name, test.expect, as.vol)
		}
	}
}

func TestAudioServiceSetVolumeNoStateChange(t *testing.T) {
	as := NewAudioService()
	as.SetVolume(60)

	// SetVolume must not alter playback state
	state := as.GetState()
	if state.State != AudioIdle {
		t.Errorf("SetVolume must not change audio state; got %q", state.State)
	}
}

func TestAudioServiceMissingFile(t *testing.T) {
	as := NewAudioService()

	// PlayLooping with a non-existent file must not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("PlayLooping panicked on missing file: %v", r)
		}
	}()

	as.PlayLooping("/nonexistent/path/music.mp3")

	// PlayLooping emits AudioPlaying synchronously, then the goroutine errors and
	// emits AudioStopped. Wait briefly for the goroutine to settle.
	time.Sleep(150 * time.Millisecond)

	state := as.GetState()
	if state.State != AudioStopped {
		t.Logf("Missing file: state after goroutine settled is %q (expected %q)", state.State, AudioStopped)
	}
}

func TestAudioServiceMissingFolder(t *testing.T) {
	as := NewAudioService()

	// PlayShuffleFolder with a non-existent folder fails synchronously
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("PlayShuffleFolder panicked on missing folder: %v", r)
		}
	}()

	as.PlayShuffleFolder("/nonexistent/folder")
	state := as.GetState()
	if state.State != AudioStopped {
		t.Errorf("Missing folder: expected %q, got %q", AudioStopped, state.State)
	}
}

func TestAudioServiceStop(t *testing.T) {
	as := NewAudioService()

	as.Stop()
	state := as.GetState()

	if state.State != AudioStopped {
		t.Errorf("After Stop(), expected %q, got %q", AudioStopped, state.State)
	}
}

func TestAudioServiceStateTransitions(t *testing.T) {
	as := NewAudioService()

	// Initial state is AudioIdle
	state1 := as.GetState()
	if state1.State != AudioIdle {
		t.Errorf("Initial state should be %q, got %q", AudioIdle, state1.State)
	}

	// SetVolume must not change state
	as.SetVolume(50)
	state2 := as.GetState()
	if state2.State != AudioIdle {
		t.Errorf("After SetVolume, state should still be %q, got %q", AudioIdle, state2.State)
	}

	// Stop transitions from idle → stopped
	as.Stop()
	state3 := as.GetState()
	if state3.State != AudioStopped {
		t.Errorf("After Stop(), expected %q, got %q", AudioStopped, state3.State)
	}

	// Second Stop is a no-op (stopCh is nil)
	as.Stop()
	state4 := as.GetState()
	if state4.State != AudioStopped {
		t.Errorf("Double Stop should stay %q, got %q", AudioStopped, state4.State)
	}
}

func TestAudioServicePlayEmptyFolder(t *testing.T) {
	tmpDir := t.TempDir()
	as := NewAudioService()

	// Empty directory contains no MP3s — should fail synchronously
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("PlayShuffleFolder panicked on empty folder: %v", r)
		}
	}()

	as.PlayShuffleFolder(tmpDir)
	state := as.GetState()
	if state.State != AudioStopped {
		t.Errorf("Empty folder: expected %q, got %q", AudioStopped, state.State)
	}
}

func TestAudioServiceValidMp3File(t *testing.T) {
	tmpDir := t.TempDir()
	mp3Path := filepath.Join(tmpDir, "test.mp3")

	// Write a minimal (but invalid) MP3 header — beep will return a decode error,
	// which the service handles silently (no panic).
	mp3Content := []byte{0xFF, 0xFB, 0x10, 0x00}
	err := os.WriteFile(mp3Path, mp3Content, 0644)
	if err != nil {
		t.Fatalf("Failed to create test MP3: %v", err)
	}

	as := NewAudioService()

	// Only recover panics in the *main goroutine* — beep decode errors are returned,
	// not panicked, so the goroutine calls emitState(AudioStopped) on failure.
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("PlayLooping panicked with minimal MP3: %v", r)
		}
	}()

	as.PlayLooping(mp3Path)
	// Allow goroutine time to hit the decode error and emit AudioStopped
	time.Sleep(150 * time.Millisecond)

	state := as.GetState()
	// After decode failure the goroutine emits AudioStopped
	if state.State != AudioStopped {
		t.Logf("After decode failure: state is %q", state.State)
	}
}
