package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSettingsServiceDefaults(t *testing.T) {
	tmpDir := t.TempDir()
	ss := &SettingsService{
		filePath: filepath.Join(tmpDir, "settings.json"),
		current:  defaultSettings(),
	}

	settings := ss.GetSettings()
	if settings.DefaultVolume != 70 {
		t.Errorf("Expected DefaultVolume 70, got %d", settings.DefaultVolume)
	}
	if !settings.AutoStartAudio {
		t.Error("AutoStartAudio should default to true")
	}
	if !settings.NotifyOnComplete {
		t.Error("NotifyOnComplete should default to true")
	}
	if settings.AutoStartNextTimer {
		t.Error("AutoStartNextTimer should default to false")
	}
	if settings.MinimizeToTray {
		t.Error("MinimizeToTray should default to false")
	}
}

func TestSettingsServiceSave(t *testing.T) {
	tmpDir := t.TempDir()
	ss := &SettingsService{
		filePath: filepath.Join(tmpDir, "settings.json"),
		current:  defaultSettings(),
	}

	modified := Settings{
		DefaultVolume:      85,
		AutoStartAudio:     false,
		NotifyOnComplete:   false,
		AutoStartNextTimer: true,
		MinimizeToTray:     true,
	}

	err := ss.SaveSettings(modified)
	if err != nil {
		t.Errorf("SaveSettings failed: %v", err)
	}

	// Verify file was created
	_, err = os.Stat(ss.filePath)
	if err != nil {
		t.Errorf("settings.json not created: %v", err)
	}

	// Verify in-memory value was updated
	current := ss.GetSettings()
	if current.DefaultVolume != 85 {
		t.Errorf("Expected DefaultVolume 85, got %d", current.DefaultVolume)
	}
	if current.AutoStartAudio {
		t.Error("AutoStartAudio should be false after save")
	}
}

func TestSettingsServiceLoad(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "settings.json")

	// Save settings with first instance
	ss1 := &SettingsService{
		filePath: filePath,
		current:  defaultSettings(),
	}
	modified := Settings{
		DefaultVolume:      90,
		AutoStartAudio:     false,
		NotifyOnComplete:   true,
		AutoStartNextTimer: false,
		MinimizeToTray:     true,
	}
	ss1.SaveSettings(modified)

	// Load with new instance
	ss2 := &SettingsService{
		filePath: filePath,
		current:  defaultSettings(),
	}
	ss2.load()

	settings := ss2.GetSettings()
	if settings.DefaultVolume != 90 {
		t.Errorf("Expected DefaultVolume 90 after load, got %d", settings.DefaultVolume)
	}
	if settings.AutoStartAudio {
		t.Error("AutoStartAudio should be false after load")
	}
	if !settings.MinimizeToTray {
		t.Error("MinimizeToTray should be true after load")
	}
}

func TestSettingsServiceVolumeRange(t *testing.T) {
	tmpDir := t.TempDir()
	ss := &SettingsService{
		filePath: filepath.Join(tmpDir, "settings.json"),
		current:  defaultSettings(),
	}

	tests := []struct {
		volume int
		valid  bool
	}{
		{0, true},
		{50, true},
		{100, true},
		{-1, true},  // JS should validate, but service stores as-is
		{150, true}, // JS should validate, but service stores as-is
	}

	for _, test := range tests {
		s := Settings{DefaultVolume: test.volume}
		err := ss.SaveSettings(s)
		if test.valid && err != nil {
			t.Errorf("Volume %d should be valid, got error: %v", test.volume, err)
		}
	}
}

func TestSettingsServicePersistence(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "settings.json")

	// Create and save
	ss1 := &SettingsService{filePath: filePath, current: defaultSettings()}
	custom := Settings{
		DefaultVolume:      75,
		AutoStartAudio:     true,
		NotifyOnComplete:   false,
		AutoStartNextTimer: true,
		MinimizeToTray:     false,
	}
	ss1.SaveSettings(custom)

	// Load in new instance
	ss2 := &SettingsService{filePath: filePath, current: defaultSettings()}
	ss2.load()
	persisted := ss2.GetSettings()

	if persisted.DefaultVolume != 75 {
		t.Errorf("Persistence: expected volume 75, got %d", persisted.DefaultVolume)
	}
	if !persisted.AutoStartNextTimer {
		t.Error("Persistence: AutoStartNextTimer should be true")
	}
}
