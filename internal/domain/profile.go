package domain

// Profile defines a named timer configuration.
type Profile struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DurationSec int    `json:"durationSec"` // total session length in seconds
	MusicPath   string `json:"musicPath"`   // file or folder path (empty = silent)
	Shuffle     bool   `json:"shuffle"`     // true = shuffle folder
}
