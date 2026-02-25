# FocusPlay Build Instructions

## Prerequisites

- **Go 1.23+** (https://go.dev/dl/)
- **Node.js 18+** (https://nodejs.org/)
- **Wails CLI** (https://wails.io/docs/gettingstarted/installation)
  - Install with: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`
- Windows: Ensure `$GOPATH/bin` (or Go install bin dir) is in your `PATH` for `wails` command

## 1. Install Go dependencies

```
go mod tidy
```

## 2. Build the frontend

```
cd frontend
npm install
npm run build
cd ..
```

## 3. Build the app (single .exe)

```
wails build
```

- Output: `build/bin/pomodoro.exe`

## 4. Run the app

```
./build/bin/pomodoro.exe
```

## 5. Run all tests (optional)

```
go test ./internal/... -v
```

---

### Troubleshooting
- If `wails` is not found, ensure Go bin dir is in your `PATH` and restart your terminal.
- If frontend build fails, check Node.js version and run `npm install` again.
- If Go build fails, run `go mod tidy` and ensure all dependencies are present.
