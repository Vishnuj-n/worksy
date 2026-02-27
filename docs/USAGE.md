# FocusPlay Usage Guide

## Getting Started

1. **Launch FocusPlay** from your Start menu or applications folder.
2. The timer defaults to 25 minutes (Pomodoro technique).
3. **Start**: Click the **Start** button or press **Space**.
4. **Pause**: Click **Pause** or press **Space** again.
5. **Stop**: Click **Stop** or press **Esc**.

---

## Features

### Profiles

Create customized timer profiles for different tasks:

1. Click the **Profiles** icon (top-left).
2. Click **+ New Profile**.
3. Set a name (e.g., "Deep Work").
4. **Duration**: Set work session length (minutes).
5. **Break**: Set break duration (minutes).
6. **Music**:
   - **File**: Loop a single MP3 track.
   - **Folder**: Shuffle songs from a folder.
   - **Break Music**: Choose separate music for breaks.
7. **Default**: Set as the default profile on launch.

### Mini Timer Mode

Keep the timer visible without distractions:

1. Click the **Mini Mode** icon (top-right) or press **M**.
2. The window shrinks to a small, always-on-top widget.
3. Drag the widget to position it anywhere on your screen.
4. Click the **Expand** icon (⛶) or press **M** again to restore the main window.

### Settings

Customize your experience via the **Settings** (gear icon):

- **Default Volume**: Set the starting volume level.
- **Auto-start Audio**: Automatically play music when the timer starts.
- **Notify on Complete**: Show a desktop notification when a session ends.
- **Auto-start Next**: Automatically begin the next session (break or work) after the current one finishes.
- **Theme**: Choose from **Dark**, **Ocean**, **Forest**, or **Minimal Black**.

---

## Keyboard Shortcuts

| Key | Action |
| :--- | :--- |
| **Space** | Start / Pause timer |
| **Esc** | Stop timer |
| **S** | Skip current session (e.g., skip break) |
| **M** | Toggle Mini Timer mode |

---

## Data & Persistence

- **Session Resume**: If you close the app mid-session, FocusPlay remembers your progress. Upon restart, a "Resume" banner appears.
- **Stats**: View your daily session count and streak at the bottom of the window.
- **Data Location**:
  - **Windows**: `%LOCALAPPDATA%\FocusPlay\state.json`
  - **macOS**: `~/Library/Caches/FocusPlay/state.json`
  - **Linux**: `~/.cache/focusplay/state.json`

---

## Troubleshooting

- **Audio not playing**: Ensure the volume slider is up and the mute button is not active. Check if the file/folder path in your profile is valid.
- **Timer resets**: If you stop the timer manually (Esc), it resets to the full duration. Pausing keeps the current time.
- **Persistence issues**: If sessions aren't saving, check permissions for the data directory listed above.
