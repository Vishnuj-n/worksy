package app

import (
	"context"

	"focusplay/internal/domain"
	"focusplay/internal/infra/events"
	"focusplay/internal/infra/storage"
	"focusplay/internal/services/audio"
	"focusplay/internal/services/persistence"
	"focusplay/internal/services/profile"
	"focusplay/internal/services/settings"
	"focusplay/internal/services/stats"
	"focusplay/internal/services/timer"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App is the root Wails application struct.
// It owns all services and exposes bound methods to the JS frontend.
type App struct {
	ctx         context.Context
	profiles    *profile.Service
	persistence *persistence.Service
	timer       *timer.Service
	audio       *audio.Service
	settings    *settings.Service
	stats       *stats.Service
}

// New creates and wires up all services.
func New() *App {
	dir := storage.DataDir()
	ps := persistence.New(dir)
	return &App{
		profiles:    profile.New(dir),
		persistence: ps,
		timer:       timer.New(ps),
		audio:       audio.New(),
		settings:    settings.New(dir),
		stats:       stats.New(dir),
	}
}

// Startup is called by Wails after the window is ready.
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	e := events.NewWailsEmitter(ctx)
	a.timer.SetEmitter(e)
	a.audio.SetEmitter(e)
	a.profiles.Load()
}

// ── Profile methods (bound to JS) ───────────────────────────────────────────

func (a *App) LoadProfiles() []domain.Profile {
	return a.profiles.Load()
}

func (a *App) SaveProfile(p domain.Profile) error {
	return a.profiles.Save(p)
}

func (a *App) GetProfileByID(id string) *domain.Profile {
	return a.profiles.GetByID(id)
}

func (a *App) DeleteProfile(id string) error {
	return a.profiles.Delete(id)
}

// ── File / folder pickers (bound to JS) ─────────────────────────────────────

func (a *App) PickMusicFile() string {
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title:   "Select MP3 file",
		Filters: []runtime.FileFilter{{DisplayName: "MP3 Audio (*.mp3)", Pattern: "*.mp3"}},
	})
	if err != nil {
		return ""
	}
	return path
}

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

func (a *App) CheckResumeSession() *domain.SessionState {
	return a.persistence.Load()
}

// ── Timer methods (bound to JS) ─────────────────────────────────────────────

func (a *App) StartTimer(profileID string, durationSec int) {
	a.timer.Start(profileID, durationSec)
}

func (a *App) ResumeTimer(state domain.SessionState) {
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

func (a *App) GetAudioState() domain.AudioStatePayload {
	return a.audio.GetState()
}

// ── Stats methods (bound to JS) ─────────────────────────────────────────────

func (a *App) GetStats() domain.StatsData {
	return a.stats.GetStats()
}

func (a *App) RecordSessionComplete() domain.StatsData {
	return a.stats.RecordSessionComplete()
}

// ── Settings methods (bound to JS) ──────────────────────────────────────────

func (a *App) GetSettings() domain.Settings {
	return a.settings.Get()
}

func (a *App) SaveSettings(s domain.Settings) error {
	return a.settings.Save(s)
}
