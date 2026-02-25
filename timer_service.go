package main

import (
	"context"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// TimerService manages the countdown timer and emits Wails events.
type TimerService struct {
	mu          sync.Mutex
	ctx         context.Context
	persistence *PersistenceService

	totalSec  int
	remainSec int
	profileID string
	running   bool
	cancel    context.CancelFunc
}

func NewTimerService(ps *PersistenceService) *TimerService {
	return &TimerService{persistence: ps}
}

// SetContext must be called from App.startup before any timer methods are used.
func (ts *TimerService) SetContext(ctx context.Context) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	ts.ctx = ctx
}

// Start begins a new countdown for durationSec seconds.
func (ts *TimerService) Start(profileID string, durationSec int) {
	ts.mu.Lock()
	if ts.cancel != nil {
		ts.cancel()
	}
	ts.profileID = profileID
	ts.totalSec = durationSec
	ts.remainSec = durationSec
	ts.running = true
	ctx, cancel := context.WithCancel(ts.ctx)
	ts.cancel = cancel
	ts.mu.Unlock()

	go ts.run(ctx)
}

// Resume restarts the timer from a previously saved state.
func (ts *TimerService) Resume(state SessionState) {
	ts.mu.Lock()
	if ts.cancel != nil {
		ts.cancel()
	}
	ts.profileID = state.ProfileID
	ts.totalSec = state.TotalSec
	ts.remainSec = state.RemainingSec
	ts.running = true
	ctx, cancel := context.WithCancel(ts.ctx)
	ts.cancel = cancel
	ts.mu.Unlock()

	go ts.run(ctx)
}

// Pause stops the tick loop while preserving remaining time.
func (ts *TimerService) Pause() {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	if ts.cancel != nil {
		ts.cancel()
		ts.cancel = nil
	}
	ts.running = false
}

// Stop halts the timer and clears persisted state.
func (ts *TimerService) Stop() {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	if ts.cancel != nil {
		ts.cancel()
		ts.cancel = nil
	}
	ts.running = false
	ts.remainSec = ts.totalSec
	ts.persistence.ClearSessionState()
}

// GetState returns current timer snapshot (safe to call from frontend).
func (ts *TimerService) GetState() map[string]interface{} {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	return map[string]interface{}{
		"running":      ts.running,
		"remainingSec": ts.remainSec,
		"totalSec":     ts.totalSec,
		"profileId":    ts.profileID,
	}
}

// ── internal ─────────────────────────────────────────────────────────────────

func (ts *TimerService) run(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	autosave := time.NewTicker(60 * time.Second)
	defer ticker.Stop()
	defer autosave.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-autosave.C:
			ts.mu.Lock()
			_ = ts.persistence.SaveSessionState(SessionState{
				ProfileID:    ts.profileID,
				TotalSec:     ts.totalSec,
				RemainingSec: ts.remainSec,
			})
			ts.mu.Unlock()

		case <-ticker.C:
			ts.mu.Lock()
			if ts.remainSec > 0 {
				ts.remainSec--
				remaining := ts.remainSec
				profileID := ts.profileID
				ts.mu.Unlock()
				runtime.EventsEmit(ts.ctx, "timerTicked", map[string]interface{}{
					"remainingSec": remaining,
					"profileId":    profileID,
				})
			} else {
				ts.running = false
				ts.mu.Unlock()
				ts.persistence.ClearSessionState()
				runtime.EventsEmit(ts.ctx, "timerCompleted", map[string]interface{}{
					"profileId": ts.profileID,
				})
				return
			}
		}
	}
}
