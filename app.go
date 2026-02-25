package main

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App is the root Wails application struct.
// It owns all services and exposes bound methods to the JS frontend.
type App struct {
	ctx         context.Context
	profiles    *ProfileManager
	persistence *PersistenceService
	timer       *TimerService
	audio       *AudioService
	settings    *SettingsService
	stats       *StatsService
}

// NewApp creates and wires up all services.
func NewApp() *App {
	ps := NewPersistenceService()
	return &App{
		profiles:    NewProfileManager(),
		persistence: ps,
		timer:       NewTimerService(ps),
		audio:       NewAudioService(),
		settings:    NewSettingsService(),
		stats:       NewStatsService(),
	}
}

// startup is called by Wails after the window is ready.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.timer.SetContext(ctx)
	a.audio.SetContext(ctx)
	// Pre-load profiles into cache
	a.profiles.LoadProfiles()
}

// ── Profile methods (bound to JS) ───────────────────────────────────────────

func (a *App) LoadProfiles() []Profile {
	return a.profiles.LoadProfiles()
}

func (a *App) SaveProfile(p Profile) error {
	return a.profiles.SaveProfile(p)
}

func (a *App) GetProfileByID(id string) *Profile {
	return a.profiles.GetProfileByID(id)
}

func (a *App) DeleteProfile(id string) error {
	return a.profiles.DeleteProfile(id)
}

// ── File / folder pickers (bound to JS) ─────────────────────────────────────

// PickMusicFile opens a native file-open dialog filtered to MP3 files.
// Returns the selected path, or "" if cancelled.
func (a *App) PickMusicFile() string {
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select MP3 file",
		Filters: []runtime.FileFilter{
			{DisplayName: "MP3 Audio (*.mp3)", Pattern: "*.mp3"},
		},
	})
	if err != nil {
		return ""
	}
	return path
}

// PickMusicFolder opens a native folder-select dialog.
// Returns the selected folder path, or "" if cancelled.
func (a *App) PickMusicFolder() string {
	path, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select music folder",
	})
	if err != nil {
		return ""
	}
	return path
}

// ── Session persistence (bound to JS) ───────────────────────────────────────

// CheckResumeSession returns a saved SessionState if one exists, or nil.
func (a *App) CheckResumeSession() *SessionState {
	return a.persistence.LoadSessionState()
}

// ── Timer methods (bound to JS) ─────────────────────────────────────────────

func (a *App) StartTimer(profileID string, durationSec int) {
	a.timer.Start(profileID, durationSec)
}

func (a *App) ResumeTimer(state SessionState) {
	a.timer.Resume(state)
}

func (a *App) PauseTimer() {
	a.timer.Pause()
}

func (a *App) StopTimer() {
	a.timer.Stop()
}

func (a *App) GetTimerState() map[string]interface{} {
	return a.timer.GetState()
}

// ── Audio methods (bound to JS) ─────────────────────────────────────────────

func (a *App) PlayLooping(filePath string) {
	a.audio.PlayLooping(filePath)
}

func (a *App) PlayShuffleFolder(folder string) {
	a.audio.PlayShuffleFolder(folder)
}

func (a *App) StopAudio() {
	a.audio.Stop()
}

func (a *App) SetVolume(v int) {
	a.audio.SetVolume(v)
}

func (a *App) GetAudioState() AudioStatePayload {
	return a.audio.GetState()
}

// ── Stats methods (bound to JS) ────────────────────────────────────────────────

// GetStats returns today's session count and current streak.
func (a *App) GetStats() StatsData {
	return a.stats.GetStats()
}

// RecordSessionComplete should be called by JS when a timer finishes naturally.
// Returns the updated stats so the frontend can update immediately.
func (a *App) RecordSessionComplete() StatsData {
	return a.stats.RecordSessionComplete()
}

// ── Settings methods (bound to JS) ──────────────────────────────────────────

func (a *App) GetSettings() Settings {
	return a.settings.GetSettings()
}

func (a *App) SaveSettings(s Settings) error {
	return a.settings.SaveSettings(s)
}
