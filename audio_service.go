package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// AudioService handles MP3 playback (single file or shuffle folder).
// Silent-fails on missing / invalid files — never crashes the app.
type AudioService struct {
	mu     sync.Mutex
	ctx    context.Context
	stopCh chan struct{}
	vol    float64 // 0.0 – 1.0
	state  AudioStatePayload
}

func NewAudioService() *AudioService {
	return &AudioService{vol: 0.7, state: AudioStatePayload{State: AudioIdle}}
}

func (a *AudioService) SetContext(ctx context.Context) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.ctx = ctx
}

// PlayLooping streams a single MP3 file in an infinite loop.
func (a *AudioService) PlayLooping(filePath string) {
	a.Stop()
	stopCh := make(chan struct{})
	a.mu.Lock()
	a.stopCh = stopCh
	a.mu.Unlock()

	trackName := filepath.Base(filePath)
	a.emitState(AudioPlaying, trackName, "Looping")

	go func() {
		for {
			select {
			case <-stopCh:
				return
			default:
				if err := a.playFile(filePath, stopCh); err != nil {
					a.emitState(AudioStopped, trackName, "Error: "+err.Error())
					return
				}
			}
		}
	}()
}

// PlayShuffleFolder scans a folder for MP3s, shuffles, plays sequentially.
func (a *AudioService) PlayShuffleFolder(folder string) {
	a.Stop()
	tracks, err := scanMP3s(folder)
	if err != nil || len(tracks) == 0 {
		a.emitState(AudioStopped, "", "No MP3s found")
		return
	}

	stopCh := make(chan struct{})
	a.mu.Lock()
	a.stopCh = stopCh
	a.mu.Unlock()

	info := fmt.Sprintf("Shuffle folder · %d tracks", len(tracks))
	shuffleStrings(tracks)
	a.emitState(AudioPlaying, filepath.Base(tracks[0]), info)

	go func() {
		idx := 0
		for {
			select {
			case <-stopCh:
				return
			default:
				track := tracks[idx%len(tracks)]
				a.emitState(AudioPlaying, filepath.Base(track), info)
				_ = a.playFile(track, stopCh)
				idx++
				if idx%len(tracks) == 0 {
					shuffleStrings(tracks)
				}
			}
		}
	}()
}

// Stop halts all playback immediately.
func (a *AudioService) Stop() {
	a.mu.Lock()
	if a.stopCh != nil {
		close(a.stopCh)
		a.stopCh = nil
	}
	a.mu.Unlock()
	speaker.Clear()
	a.emitState(AudioStopped, "", "")
}

// SetVolume adjusts playback volume (0–100 int from frontend).
func (a *AudioService) SetVolume(v int) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if v < 0 {
		v = 0
	}
	if v > 100 {
		v = 100
	}
	a.vol = float64(v) / 100.0
}

// GetState returns the current audio state for the frontend.
func (a *AudioService) GetState() AudioStatePayload {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.state
}

// ── internal ─────────────────────────────────────────────────────────────────

func (a *AudioService) playFile(path string, stopCh chan struct{}) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		return err
	}
	defer streamer.Close()

	_ = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	done := make(chan struct{})
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		close(done)
	})))

	select {
	case <-done:
		return nil
	case <-stopCh:
		speaker.Clear()
		return nil
	}
}

func (a *AudioService) emitState(state AudioPlaybackState, track, info string) {
	a.mu.Lock()
	a.state = AudioStatePayload{State: state, TrackName: track, TrackInfo: info}
	ctx := a.ctx
	payload := a.state
	a.mu.Unlock()

	if ctx != nil {
		runtime.EventsEmit(ctx, "audioStateChanged", payload)
	}
}

func scanMP3s(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var files []string
	for _, e := range entries {
		if !e.IsDir() && strings.ToLower(filepath.Ext(e.Name())) == ".mp3" {
			files = append(files, filepath.Join(dir, e.Name()))
		}
	}
	return files, nil
}

func shuffleStrings(s []string) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(s), func(i, j int) { s[i], s[j] = s[j], s[i] })
}
