package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestPersistenceServiceSaveSessionState(t *testing.T) {
	tmpDir := t.TempDir()
	ps := &PersistenceService{filePath: filepath.Join(tmpDir, "state.json")}

	state := SessionState{
		ProfileID:    "pomodoro",
		TotalSec:     25 * 60,
		RemainingSec: 1200,
	}

	err := ps.SaveSessionState(state)
	if err != nil {
		t.Errorf("SaveSessionState failed: %v", err)
	}

	// Verify file was created
	_, err = os.Stat(ps.filePath)
	if err != nil {
		t.Errorf("state.json not created: %v", err)
	}
}

func TestPersistenceServiceLoadSessionState(t *testing.T) {
	tmpDir := t.TempDir()
	ps := &PersistenceService{filePath: filepath.Join(tmpDir, "state.json")}

	original := SessionState{
		ProfileID:    "pomodoro",
		TotalSec:     1500,
		RemainingSec: 1200,
	}
	ps.SaveSessionState(original)

	loaded := ps.LoadSessionState()
	if loaded == nil {
		t.Error("Failed to load saved session state")
	}
	if loaded.ProfileID != "pomodoro" {
		t.Errorf("Expected ProfileID 'pomodoro', got %q", loaded.ProfileID)
	}
	if loaded.RemainingSec != 1200 {
		t.Errorf("Expected 1200 remaining, got %d", loaded.RemainingSec)
	}
}

func TestPersistenceServiceStaleSessionExpiry(t *testing.T) {
	tmpDir := t.TempDir()
	ps := &PersistenceService{filePath: filepath.Join(tmpDir, "state.json")}

	// Create a session with stale timestamp (>24h old)
	state := SessionState{
		ProfileID:    "pomodoro",
		TotalSec:     1500,
		RemainingSec: 1200,
		SavedAt:      time.Now().Unix() - 86401, // 24h + 1 second
	}
	ps.SaveSessionState(state)

	// Loading should return nil and delete the file
	loaded := ps.LoadSessionState()
	if loaded != nil {
		t.Error("Expected nil for stale session")
	}

	// File should be deleted
	_, err := os.Stat(ps.filePath)
	if !os.IsNotExist(err) {
		t.Error("Stale state.json was not deleted")
	}
}

func TestPersistenceServiceFreshSession(t *testing.T) {
	tmpDir := t.TempDir()
	ps := &PersistenceService{filePath: filepath.Join(tmpDir, "state.json")}

	// Create a fresh session (within 24h)
	state := SessionState{
		ProfileID:    "pomodoro",
		TotalSec:     1500,
		RemainingSec: 1200,
		SavedAt:      time.Now().Unix() - 3600, // 1 hour old
	}
	ps.SaveSessionState(state)

	loaded := ps.LoadSessionState()
	if loaded == nil {
		t.Error("Fresh session should be loaded")
	}
	if loaded.RemainingSec != 1200 {
		t.Errorf("Expected 1200 remaining, got %d", loaded.RemainingSec)
	}
}

func TestPersistenceServiceClearSessionState(t *testing.T) {
	tmpDir := t.TempDir()
	ps := &PersistenceService{filePath: filepath.Join(tmpDir, "state.json")}

	state := SessionState{ProfileID: "test", TotalSec: 100, RemainingSec: 50}
	ps.SaveSessionState(state)

	// Verify file exists
	_, err := os.Stat(ps.filePath)
	if err != nil {
		t.Error("state.json not created before clear")
	}

	// Clear
	ps.ClearSessionState()

	// File should be deleted
	_, err = os.Stat(ps.filePath)
	if !os.IsNotExist(err) {
		t.Error("state.json not deleted after ClearSessionState")
	}
}

func TestPersistenceServiceNoFileFallback(t *testing.T) {
	tmpDir := t.TempDir()
	ps := &PersistenceService{filePath: filepath.Join(tmpDir, "nonexistent.json")}

	// Loading when file doesn't exist should return nil
	loaded := ps.LoadSessionState()
	if loaded != nil {
		t.Error("Expected nil when state file doesn't exist")
	}
}
