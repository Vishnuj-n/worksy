# Plan: FocusPlay MVP Service Architecture (Go + Wails 2) - Instructions Only

**TL;DR:** Build FocusPlay as Wails 2 app (Go backend + HTML/CSS/JS frontend). Four services communicate via Wails events. JSON only in ProfileManager/PersistenceService. 20MB native Windows executable.

***

## Steps

### Phase 1: Project Setup

1. Initialize Wails project with vanilla template
2. Add Go audio dependencies (go-mp3, beep, flock)
3. Create AppData directories for profiles.json and state.json
4. Set up project structure: app.go, models.go, services.go, frontend/

### Phase 2: Go Services Layer

5. Define models: Profile struct, SessionState struct, AudioPlaybackState enum

6. Implement ProfileManager:
   - LoadProfiles(): Read profiles.json from AppData, return cached list
   - SaveProfile(): Write single profile to profiles.json
   - GetProfileById(): Return from memory cache only

7. Implement PersistenceService:
   - LoadSessionState(): Read state.json from AppData
   - SaveSessionState(): Write state.json every 60 seconds
   - ClearSessionState(): Delete state.json on completion

8. Implement TimerService:
   - Start(duration): Begin countdown with time.Timer
   - Pause/Resume/Stop controls
   - Emit "timerTicked" event every second with remaining time
   - Emit "timerCompleted" event at zero
   - Auto-save via PersistenceService every 60s

9. Implement AudioService:
   - PlayLooping(filePath): Stream MP3 with infinite loop
   - PlayShuffleFolder(folder): Scan MP3s, shuffle, play sequentially
   - Stop(): Halt playback immediately
   - SetVolume(): Control playback volume
   - Emit "audioStateChanged" events
   - Silent fail on missing/invalid files

### Phase 3: Frontend (HTML/CSS/JS)

10. Create glassmorphism UI:
    - Large centered timer display (HH:MM:SS)
    - Profile dropdown selector
    - Start/Pause/Stop buttons
    - Audio status indicator
    - Mica backdrop blur effects
    - Segoe UI font + Fluent spacing

11. Frontend bindings:
    - Call Go methods: LoadProfiles(), StartTimer(), PlayAudio()
    - Listen to events: timerTicked, timerCompleted, audioStateChanged
    - Format remaining seconds as MM:SS display
    - Update UI reactively on events

### Phase 4: App Integration

12. Main App struct lifecycle:
    - WailsInit(): Load profiles, check resume session, start event loop
    - Startup: Show resume prompt if valid session state exists
    - Expose all services as bindable methods to frontend
    - Handle all Wails events and method calls

### Phase 5: Build & Deploy

13. Development workflow:
    - wails dev (hot reload frontend + Go backend)
    - Real-time UI updates during development

14. Production build:
    - wails build (single 20MB Windows executable)
    - Self-contained, no runtime dependencies

***

## Verification Checklist

- Timer countdown works + persists across app restarts
- Audio plays single MP3s and shuffled folders correctly
- Profiles load/save from AppData JSON files
- Session resume prompt appears on startup
- Missing music files handled silently (no crash)
- Glassmorphism UI responsive + modern appearance
- Single executable runs on clean Windows 11

***

## Key Decisions

- Wails 2: Native window frame + hot reload + Go performance
- HTML/CSS/JS frontend: Faster modern styling than XAML
- Event-driven: Go emits → JS subscribes (same architecture as C# plan)
- JSON encapsulation preserved: Only two services touch JSON
- 20MB binary: Smaller than .NET NativeAOT

***

**Result:** Native Windows app with modern glassmorphism UI, full FocusPlay feature parity, production-ready in one executable.

**Commands:** `wails dev` → `wails build` 🚀