
---

## Overall Status

| Category | Count | Status |
|----------|-------|--------|
| **Core Features (18)** | 18/18 | ✅ Complete |
| **P0 Bugs** | 2/2 | ✅ Fixed |
| **P1 Partial Wiring** | 3/4 | ✅ Done (1 Removed) |
| **P2 New Features** | 2/4 | 50% Done |
| **P3 Nice-to-Have** | — | Not started |

---

| # | Feature | Location |
|---|---------|----------|
| 1 | Timer countdown (start/pause/stop/resume) | `internal/services/timer/` |
| 2 | Wails event system (timerTicked, timerCompleted, audioStateChanged) | `internal/infra/events/` |
| 3 | Profile CRUD (create, edit, delete, 3 defaults) | `internal/services/profile/` |
| 4 | Profile dropdown + switching in UI | `frontend/src/main.js` |
| 5 | Audio: single-file infinite loop | `internal/services/audio/` |
| 6 | Audio: shuffle folder playback | `internal/services/audio/` |
| 7 | Volume control (slider + settings default) | `audio/SetVolume` + frontend |
| 8 | Session persistence (auto-save every 60s, 24h expiry) | `internal/services/persistence/` |
| 9 | Session resume on startup (banner + Resume button) | `frontend/src/main.js init()` |
| 10 | Settings panel (4 toggles + volume slider, persisted) | `internal/services/settings/` |
| 11 | Stats tracking (sessions today, streak, day rollover) | `internal/services/stats/` |
| 12 | Glassmorphism UI (animated blobs, blur, glass card) | `frontend/src/style.css` |
| 13 | Native file/folder picker dialogs | `app.go PickMusicFile/PickMusicFolder` |
| 14 | Silent-fail audio (missing/invalid files don't crash) | `audio/service.go` |
| 15 | NSIS Windows installer config | `build/windows/installer/project.nsi` |
| 16 | Unit tests for all 6 services (32 tests pass) | `internal/services/*/service_test.go` |
| 17 | JSON encapsulation (only storage layer touches JSON) | `internal/infra/storage/` |
| 18 | Clean package structure (domain → services → app) | `internal/` |

---

## Bugs to Fix — CLOSED ✅

### ✅ B1. Volume slider has no audible effect
- **FIXED:** Wrapped streamer with `effects.Volume` in `audio/service.go playFile()`. Updates applied live via `speaker.Lock()`.

### ✅ B2. OS notifications may silently fail
- **FIXED:** Added `Notification.requestPermission()` in `init()`.

---

## Partially Wired (P1) — MOSTLY COMPLETE ✅

### ❌ P1.1 Wire minimize-to-tray
- **REMOVED:** Feature removed in favor of global hotkey/standard window behavior. Code deleted from backend & frontend.

### ✅ P1.2 Resume audio on session resume
- **FIXED:** Resume button now looks up profile's music path and calls `PlayLooping`/`PlayShuffleFolder`.

### ✅ P1.3 Stop audio on timer pause
- **FIXED:** Pause handler now calls `StopAudio()` alongside `PauseTimer()`.

### P1.4 Auto-start next session has no break interval
- **Status:** `autoStartNextTimer` immediately restarts the same profile — no break concept.
- **Fix:** Defer to P3 (work/break cycle) or add a simple 5-minute cooldown.

---

## New Features — In Progress (P2)

### ✅ 6. Keyboard shortcuts
- **DONE:** Space = start/pause toggle, Escape = stop timer, S = skip. Implemented in `main.js` with panel focus check.

### 7. Custom app icon
- Create `build/appicon.png` (1024×1024) and `build/windows/icon.ico`
- Update `wails.json` with icon reference

### ✅ 8. Fix wails.json metadata
- **DONE:** Added `info` block with `ProductName`, `ProductVersion`, `Copyright`, `Comments`.

### 9. Completion chime
- Embed a short notification sound (small MP3, ~50KB) in the Go binary
- Play it via beep when `timerCompleted` fires, independent of background music
- Respect a new `playSoundOnComplete` setting

---

## Nice-to-Have (P3)

| # | Feature | Notes |
|---|---------|-------|
| 10 | Dark/light theme | Extract colors to CSS custom properties, add toggle, respect `prefers-color-scheme` |
| 11 | Work/break cycle | Add `breakDurationSec` to Profile, auto-start break after work, round counter |
| 12 | Taskbar progress bar | Needs CGo call to `ITaskbarList3::SetProgressValue` — Wails doesn't expose this |
| 13 | Stats visualization | Weekly chart, calendar heatmap, CSV export |
| 14 | Frameless window | Set `Frameless: true`, add `--wails-draggable` header, custom close/minimize buttons |
| 15 | Multi-platform | Test on macOS/Linux (beep library supports all three) |

---

## Verification

```bash
# Run all tests
go test ./internal/... -v

# Development with hot reload
wails dev

# Production build (single .exe)
wails build
# Output: build/bin/POMODORO.exe