package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

// ProfileManager handles loading and saving profiles from AppData/profiles.json.
// JSON is ONLY touched here — all other services use in-memory Profile structs.
type ProfileManager struct {
	mu       sync.RWMutex
	profiles []Profile
	filePath string
}

func NewProfileManager() *ProfileManager {
	appData, _ := os.UserCacheDir() // %LOCALAPPDATA%
	dir := filepath.Join(appData, "FocusPlay")
	_ = os.MkdirAll(dir, 0755)
	return &ProfileManager{
		filePath: filepath.Join(dir, "profiles.json"),
	}
}

// LoadProfiles reads profiles.json and caches the result. Returns defaults on first run.
func (pm *ProfileManager) LoadProfiles() []Profile {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	data, err := os.ReadFile(pm.filePath)
	if err != nil {
		pm.profiles = defaultProfiles()
		pm.saveUnlocked()
		return pm.profiles
	}
	if err := json.Unmarshal(data, &pm.profiles); err != nil || len(pm.profiles) == 0 {
		pm.profiles = defaultProfiles()
	}
	return pm.profiles
}

// SaveProfile upserts a profile in the cache and writes profiles.json.
func (pm *ProfileManager) SaveProfile(p Profile) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	for i, existing := range pm.profiles {
		if existing.ID == p.ID {
			pm.profiles[i] = p
			return pm.saveUnlocked()
		}
	}
	pm.profiles = append(pm.profiles, p)
	return pm.saveUnlocked()
}

// GetProfileByID returns a profile from the in-memory cache only.
func (pm *ProfileManager) GetProfileByID(id string) *Profile {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	for _, p := range pm.profiles {
		if p.ID == id {
			cp := p
			return &cp
		}
	}
	return nil
}

// DeleteProfile removes a profile by ID and persists the change.
func (pm *ProfileManager) DeleteProfile(id string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	filtered := pm.profiles[:0]
	for _, p := range pm.profiles {
		if p.ID != id {
			filtered = append(filtered, p)
		}
	}
	pm.profiles = filtered
	return pm.saveUnlocked()
}

func (pm *ProfileManager) saveUnlocked() error {
	data, err := json.MarshalIndent(pm.profiles, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(pm.filePath, data, 0644)
}

func defaultProfiles() []Profile {
	return []Profile{
		{ID: "deep-work", Name: "Deep Work — 90 min", DurationSec: 90 * 60},
		{ID: "pomodoro", Name: "Pomodoro — 25 min", DurationSec: 25 * 60},
		{ID: "short-break", Name: "Short Break — 5 min", DurationSec: 5 * 60},
	}
}
