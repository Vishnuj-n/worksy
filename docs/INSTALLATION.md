# FocusPlay Build & Installation Guide

## Prerequisites

- **Go 1.23+** (https://go.dev/dl/)
- **Node.js 18+** (https://nodejs.org/)
- **Wails CLI**: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`

Ensure the Go binary directory (e.g., `$GOPATH/bin`) is in your system `PATH`.

---

## 1. Install Dependencies

In the root of the project:

```bash
# Install Go dependencies
go mod tidy

# Install Frontend dependencies
cd frontend
npm install
cd ..
```

---

## 2. Build the Application

### Windows (Recommended)

Generate a single `.exe` file and an NSIS installer:

```bash
wails build --nsis
```

- **Output:** `build/bin/focusplay.exe` and the installer (e.g., `focusplay-amd64-installer.exe`).

### macOS

Build a universal binary or `.app` bundle:

```bash
wails build --platform darwin/universal
```

### Linux

Build a binary for AMD64:

```bash
wails build --platform linux/amd64
```

---

## 3. Run Development Mode

To run the application with hot-reload for frontend changes:

```bash
wails dev
```

This starts a Vite server and rebuilds the Go backend automatically on changes.

---

## Troubleshooting

### `wails` command not found
- Ensure your Go binary path (e.g., `~/go/bin` or `%GOPATH%\bin`) is added to your system `PATH`.
- Restart your terminal or computer after installing Wails.

### Frontend build fails
- Check your Node.js version (`node -v`). FocusPlay requires Node.js 18 or later.
- Delete `frontend/node_modules` and run `npm install` again.

### Go build fails
- Run `go mod tidy` to ensure all dependencies are downloaded.
- Ensure you have a C compiler installed (e.g., `gcc` on Linux/macOS, TDM-GCC or MinGW on Windows) if required by dependencies (though Wails often handles this).

### Application crashes on start
- Check the console output if running via `wails dev`.
- Ensure you have write permissions to the application data directory (usually `%LOCALAPPDATA%\FocusPlay` on Windows, or `~/.cache/focusplay` on Linux/macOS).
