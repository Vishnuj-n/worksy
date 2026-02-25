package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// PersistenceService handles reading and writing state.json.
// JSON is ONLY touched here — all other services use in-memory SessionState.
type PersistenceService struct {
	mu       sync.Mutex
	filePath string
}

func NewPersistenceService() *PersistenceService {
	appData, _ := os.UserCacheDir()
	dir := filepath.Join(appData, "FocusPlay")
	_ = os.MkdirAll(dir, 0755)
	return &PersistenceService{
		filePath: filepath.Join(dir, "state.json"),
	}
}

// LoadSessionState reads state.json. Returns nil if file is absent or stale (>24h).
func (ps *PersistenceService) LoadSessionState() *SessionState {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	data, err := os.ReadFile(ps.filePath)
	if err != nil {
		return nil
	}
	var s SessionState
	if err := json.Unmarshal(data, &s); err != nil {
		return nil
	}
	// Discard sessions saved more than 24 hours ago
	if s.SavedAt == 0 || time.Now().Unix()-s.SavedAt > 86400 {
		_ = os.Remove(ps.filePath)
		return nil
	}
	return &s
}

// SaveSessionState writes state.json with the current Unix timestamp.
func (ps *PersistenceService) SaveSessionState(s SessionState) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	s.SavedAt = time.Now().Unix()
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(ps.filePath, data, 0644)
}

// ClearSessionState deletes state.json (called on session completion).
func (ps *PersistenceService) ClearSessionState() {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	_ = os.Remove(ps.filePath)
}
