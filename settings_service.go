package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

// Settings holds global app preferences persisted to settings.json.
type Settings struct {
	DefaultVolume      int  `json:"defaultVolume"`      // 0-100
	AutoStartAudio     bool `json:"autoStartAudio"`     // play audio automatically when timer starts
	NotifyOnComplete   bool `json:"notifyOnComplete"`   // OS notification on session complete
	AutoStartNextTimer bool `json:"autoStartNextTimer"` // immediately start next session
	MinimizeToTray     bool `json:"minimizeToTray"`     // minimize to system tray on close
}

// SettingsService persists and exposes Settings.
// JSON is ONLY touched here.
type SettingsService struct {
	mu       sync.RWMutex
	current  Settings
	filePath string
}

func NewSettingsService() *SettingsService {
	appData, _ := os.UserCacheDir()
	dir := filepath.Join(appData, "FocusPlay")
	_ = os.MkdirAll(dir, 0755)
	ss := &SettingsService{
		filePath: filepath.Join(dir, "settings.json"),
		current:  defaultSettings(),
	}
	ss.load()
	return ss
}

func (ss *SettingsService) GetSettings() Settings {
	ss.mu.RLock()
	defer ss.mu.RUnlock()
	return ss.current
}

func (ss *SettingsService) SaveSettings(s Settings) error {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	ss.current = s
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(ss.filePath, data, 0644)
}

func (ss *SettingsService) load() {
	data, err := os.ReadFile(ss.filePath)
	if err != nil {
		return
	}
	_ = json.Unmarshal(data, &ss.current)
}

func defaultSettings() Settings {
	return Settings{
		DefaultVolume:      70,
		AutoStartAudio:     true,
		NotifyOnComplete:   true,
		AutoStartNextTimer: false,
		MinimizeToTray:     false,
	}
}
