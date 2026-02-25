# FocusPlay — Project Status & Roadmap

> Wails 2 (Go 1.23 backend + Vite/HTML/CSS/JS frontend)  
> Single 20 MB Windows executable · 32/32 tests passing  
> `wails dev` → hot reload · `wails build` → `build/bin/POMODORO.exe`

---

## Architecture

```
focusplay/
├── main.go                        Wails bootstrap
├── internal/
│   ├── domain/                    Pure types (no deps)
│   │   ├── profile.go             Profile struct
│   │   ├── session.go             SessionState, StatsData
│   │   ├── audio.go               AudioPlaybackState, AudioStatePayload
│   │   └── settings.go            Settings + DefaultSettings()
│   ├── infra/
│   │   ├── storage/json_store.go  Load / Save / DataDir (%LOCALAPPDATA%\FocusPlay\)
│   │   └── events/emitter.go      Emitter interface, WailsEmitter, Noop
│   ├── services/
│   │   ├── profile/               CRUD for profiles.json
│   │   ├── persistence/           state.json (auto-save + 24 h expiry)
│   │   ├── timer/                 countdown, Pause/Resume/Stop, timerTicked/timerCompleted events
│   │   ├── audio/                 MP3 loop + shuffle-folder, effects.Volume real-time control
│   │   ├── stats/                 sessions-today, rolling streak, stats.json
│   │   └── settings/              settings.json defaults + persistence
│   └── app/app.go                 App struct — 40+ Wails-bound methods, BeforeClose hook
└── frontend/
    ├── index.html                 Main card, profile panel, settings panel, resume banner
    ├── src/style.css              Glassmorphism, blobs, toggles (~700 lines)
    └── src/main.js                Event-driven JS, keyboard shortcuts (~430 lines)
```

---

## Completed Features

- [x] Glassmorphism UI (animated blobs, backdrop blur, Segoe UI)
- [x] Profile management panel (CRUD, file/folder picker dialogs)
- [x] Timer countdown with MM:SS display + progress bar
- [x] Pause / Resume / Stop / Skip controls
- [x] Session persistence — auto-save every 60 s, 24 h expiry, resume banner on restart
- [x] MP3 looping + shuffle-folder playback (gopxl/beep)
- [x] Volume slider — **real-time** via `effects.Volume` + `speaker.Lock()` (fixed B1)
- [x] Daily stats (sessions today, streak) in footer
- [x] Settings panel (5 prefs: volume, auto-audio, notifications, auto-next, minimize-to-tray)
- [x] OS desktop notifications on session complete (with `requestPermission()` on boot — fixed B2)
- [x] Minimize-to-tray (`OnBeforeClose` hides window when setting enabled)
- [x] Audio restarts when resuming a saved session (P1.2)
- [x] Audio stops when timer is paused (P1.3)
- [x] Keyboard shortcuts: `Space` = start/pause, `Esc` = stop, `S` = skip
- [x] Auto-start next timer setting wired
- [x] wails.json product metadata block (name, version, copyright)
- [x] 32 unit tests across all 6 services

---

## Possible Future Enhancements (P3)

- [ ] Custom app icon (replace default Wails icon)
- [ ] Completion chime (short beep on `timerCompleted`)
- [ ] Dark/light theme toggle (CSS custom properties)
- [ ] System tray icon with right-click menu (needs third-party tray library)
- [ ] Pomodoro cycle tracking (work / short break / long break phases)
- [ ] Export stats to CSV

---

## Commands

```sh
# Development (hot reload)
wails dev

# Production build → build/bin/POMODORO.exe
wails build

# Tests
go test ./internal/... -v

# Full build check
go build ./...
```