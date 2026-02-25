package domain

// Settings holds global app preferences persisted to settings.json.
type Settings struct {
	DefaultVolume      int  `json:"defaultVolume"` // 0-100
	AutoStartAudio     bool `json:"autoStartAudio"`
	NotifyOnComplete   bool `json:"notifyOnComplete"`
	AutoStartNextTimer bool `json:"autoStartNextTimer"`
	MinimizeToTray     bool `json:"minimizeToTray"`
}

// DefaultSettings returns the factory defaults shown on first run.
func DefaultSettings() Settings {
	return Settings{
		DefaultVolume:      70,
		AutoStartAudio:     true,
		NotifyOnComplete:   true,
		AutoStartNextTimer: false,
		MinimizeToTray:     false,
	}
}
