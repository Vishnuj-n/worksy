package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProfileManagerLoadProfiles(t *testing.T) {
	tmpDir := t.TempDir()
	pm := &ProfileManager{filePath: filepath.Join(tmpDir, "profiles.json")}

	// Load returns default profiles when file doesn't exist
	profiles := pm.LoadProfiles()
	if len(profiles) != 3 {
		t.Errorf("Expected 3 default profiles, got %d", len(profiles))
	}
	if profiles[0].ID != "deep-work" {
		t.Errorf("Expected first profile ID 'deep-work', got %q", profiles[0].ID)
	}
	if profiles[0].DurationSec != 90*60 {
		t.Errorf("Expected 5400 sec, got %d", profiles[0].DurationSec)
	}

	// Verify profiles.json was created
	_, err := os.Stat(pm.filePath)
	if err != nil {
		t.Errorf("profiles.json not created: %v", err)
	}
}

func TestProfileManagerSaveProfile(t *testing.T) {
	tmpDir := t.TempDir()
	pm := &ProfileManager{filePath: filepath.Join(tmpDir, "profiles.json")}
	pm.LoadProfiles() // Initialize defaults

	newProfile := Profile{
		ID:          "custom",
		Name:        "Custom Session",
		DurationSec: 45 * 60,
		MusicPath:   "/path/to/music.mp3",
		Shuffle:     false,
	}

	err := pm.SaveProfile(newProfile)
	if err != nil {
		t.Errorf("SaveProfile failed: %v", err)
	}

	// Verify it was added to cache
	found := pm.GetProfileByID("custom")
	if found == nil {
		t.Error("Profile not found after save")
	}
	if found.Name != "Custom Session" {
		t.Errorf("Expected name 'Custom Session', got %q", found.Name)
	}
}

func TestProfileManagerUpdateProfile(t *testing.T) {
	tmpDir := t.TempDir()
	pm := &ProfileManager{filePath: filepath.Join(tmpDir, "profiles.json")}
	pm.LoadProfiles()

	// Update existing profile
	updated := Profile{
		ID:          "pomodoro",
		Name:        "Pomodoro — 35 min",
		DurationSec: 35 * 60,
	}
	pm.SaveProfile(updated)

	retrieved := pm.GetProfileByID("pomodoro")
	if retrieved.DurationSec != 35*60 {
		t.Errorf("Expected 2100 sec, got %d", retrieved.DurationSec)
	}
}

func TestProfileManagerDeleteProfile(t *testing.T) {
	tmpDir := t.TempDir()
	pm := &ProfileManager{filePath: filepath.Join(tmpDir, "profiles.json")}
	pm.LoadProfiles()

	initialCount := len(pm.profiles)
	err := pm.DeleteProfile("pomodoro")
	if err != nil {
		t.Errorf("DeleteProfile failed: %v", err)
	}

	if len(pm.profiles) != initialCount-1 {
		t.Errorf("Expected %d profiles after delete, got %d", initialCount-1, len(pm.profiles))
	}

	found := pm.GetProfileByID("pomodoro")
	if found != nil {
		t.Error("Deleted profile still exists")
	}
}

func TestProfileManagerGetProfileByID(t *testing.T) {
	tmpDir := t.TempDir()
	pm := &ProfileManager{filePath: filepath.Join(tmpDir, "profiles.json")}
	pm.LoadProfiles()

	profile := pm.GetProfileByID("deep-work")
	if profile == nil {
		t.Error("Failed to get existing profile")
	}
	if profile.ID != "deep-work" {
		t.Errorf("Expected ID 'deep-work', got %q", profile.ID)
	}

	// Non-existent profile
	notFound := pm.GetProfileByID("nonexistent")
	if notFound != nil {
		t.Error("Expected nil for non-existent profile")
	}
}

func TestProfileManagerPersistence(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "profiles.json")

	// Create and save a profile
	pm1 := &ProfileManager{filePath: filePath}
	pm1.LoadProfiles()
	custom := Profile{ID: "test", Name: "Test", DurationSec: 120}
	pm1.SaveProfile(custom)

	// Create new manager instance from same file
	pm2 := &ProfileManager{filePath: filePath}
	profiles := pm2.LoadProfiles()

	// Should load previously saved profile
	found := false
	for _, p := range profiles {
		if p.ID == "test" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Persistence failed: saved profile not loaded")
	}
}
