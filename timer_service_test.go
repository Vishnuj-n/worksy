package main

import (
	"context"
	"testing"
	"time"
)

func TestTimerServiceStart(t *testing.T) {
	ps := NewPersistenceService()
	ts := NewTimerService(ps)
	ctx := context.Background()
	ts.SetContext(ctx)

	ts.Start("test-profile", 60)

	state := ts.GetState()
	if !state["running"].(bool) {
		t.Error("Timer should be running after Start")
	}
	if state["totalSec"].(int) != 60 {
		t.Errorf("Expected totalSec 60, got %v", state["totalSec"])
	}
	if state["remainingSec"].(int) != 60 {
		t.Errorf("Expected remainingSec 60 at start, got %v", state["remainingSec"])
	}
}

func TestTimerServicePause(t *testing.T) {
	ps := NewPersistenceService()
	ts := NewTimerService(ps)
	ctx := context.Background()
	ts.SetContext(ctx)

	ts.Start("test-profile", 60)
	time.Sleep(100 * time.Millisecond)

	ts.Pause()
	state := ts.GetState()
	if state["running"].(bool) {
		t.Error("Timer should not be running after Pause")
	}

	remainingBeforePause := state["remainingSec"].(int)

	// Wait a bit and verify time didn't decrease
	time.Sleep(100 * time.Millisecond)
	stateLater := ts.GetState()
	if stateLater["remainingSec"].(int) != remainingBeforePause {
		t.Error("Time should not decrease after pause")
	}
}

func TestTimerServiceStop(t *testing.T) {
	ps := NewPersistenceService()
	ts := NewTimerService(ps)
	ctx := context.Background()
	ts.SetContext(ctx)

	ts.Start("test-profile", 60)

	ts.Stop()
	state := ts.GetState()
	if state["running"].(bool) {
		t.Error("Timer should not be running after Stop")
	}
	if state["remainingSec"].(int) != 60 {
		t.Errorf("After stop, remaining should reset to 60, got %v", state["remainingSec"])
	}
}

func TestTimerServiceGetState(t *testing.T) {
	ps := NewPersistenceService()
	ts := NewTimerService(ps)
	ctx := context.Background()
	ts.SetContext(ctx)

	ts.Start("test-profile", 120)

	state := ts.GetState()
	if state["profileId"].(string) != "test-profile" {
		t.Errorf("Expected profileId 'test-profile', got %v", state["profileId"])
	}
	if state["totalSec"].(int) != 120 {
		t.Errorf("Expected totalSec 120, got %v", state["totalSec"])
	}
}

func TestTimerServiceResume(t *testing.T) {
	ps := NewPersistenceService()
	ts := NewTimerService(ps)
	ctx := context.Background()
	ts.SetContext(ctx)

	sessionState := SessionState{
		ProfileID:    "pomodoro",
		TotalSec:     1500,
		RemainingSec: 1200,
	}

	ts.Resume(sessionState)

	state := ts.GetState()
	if !state["running"].(bool) {
		t.Error("Timer should be running after Resume")
	}
	if state["totalSec"].(int) != 1500 {
		t.Errorf("Expected totalSec 1500 after resume, got %v", state["totalSec"])
	}
	if state["remainingSec"].(int) != 1200 {
		t.Errorf("Expected remainingSec 1200 after resume, got %v", state["remainingSec"])
	}
}

func TestTimerServiceMultipleStarts(t *testing.T) {
	ps := NewPersistenceService()
	ts := NewTimerService(ps)
	ctx := context.Background()
	ts.SetContext(ctx)

	// Start first timer
	ts.Start("profile1", 60)
	state1 := ts.GetState()

	// Start second timer (should cancel first)
	ts.Start("profile2", 120)
	state2 := ts.GetState()

	if state2["profileId"].(string) != "profile2" {
		t.Error("Second start should override first")
	}
	if state2["totalSec"].(int) != 120 {
		t.Errorf("Expected new duration 120, got %v", state2["totalSec"])
	}

	_ = state1 // avoid unused variable
}

func TestTimerServiceNilContext(t *testing.T) {
	ps := NewPersistenceService()
	ts := NewTimerService(ps)

	// Should not panic even without SetContext being called
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Timer panicked without context: %v", r)
		}
	}()

	ts.Start("test", 60)
}
