package main

// ── Profile ─────────────────────────────────────────────────────────────────

// Profile defines a named timer configuration.
type Profile struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DurationSec int    `json:"durationSec"` // total session length in seconds
	MusicPath   string `json:"musicPath"`   // file or folder path (empty = silent)
	Shuffle     bool   `json:"shuffle"`     // true = shuffle folder
}

// ── SessionState ─────────────────────────────────────────────────────────────

// SessionState is persisted to state.json so the session survives restarts.
// SavedAt is a Unix timestamp (int64) to avoid Wails binding issues with time.Time.
type SessionState struct {
	ProfileID    string `json:"profileId"`
	TotalSec     int    `json:"totalSec"`
	RemainingSec int    `json:"remainingSec"`
	SavedAt      int64  `json:"savedAt"`
}

// ── AudioPlaybackState ───────────────────────────────────────────────────────

type AudioPlaybackState string

const (
	AudioIdle    AudioPlaybackState = "idle"
	AudioPlaying AudioPlaybackState = "playing"
	AudioStopped AudioPlaybackState = "stopped"
)

// AudioStatePayload is emitted via the "audioStateChanged" event.
type AudioStatePayload struct {
	State     AudioPlaybackState `json:"state"`
	TrackName string             `json:"trackName"`
	TrackInfo string             `json:"trackInfo"` // e.g. "Shuffle folder · 12 tracks"
}
