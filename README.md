# FocusPlay — Pomodoro Timer

**v1.0.0** | A lightweight Pomodoro timer application built with Wails and Go.

## About

FocusPlay is a distraction-free Pomodoro timer that helps you manage work sessions with customizable breaks and background music/playlists. The app persists your session history and supports multiple user profiles with independent preferences.

## Features

- **Pomodoro workflow**: Start, pause, stop sessions with customizable work and break durations
- **Multiple profiles**: Create and manage independent user profiles with separate settings
- **Session persistence**: Automatically saves progress; resume your work across app restarts
- **Background music/playlists**: Play looping single tracks or shuffle folder playlists during work sessions
- **Statistics**: Track session counts and productivity metrics
- **Windows installer**: Built with NSIS for easy installation and cleanup

## Quick Start

1. Download the Windows installer from `build/windows/installer/` or the latest release
2. Run the installer to install FocusPlay
3. Launch the app and start your first Pomodoro session

Alternatively, to build from source:
- Go 1.19 or later
- Node.js 16+
- Wails CLI: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`

You can configure the project by editing `wails.json`. More information about the project settings can be found
here: https://wails.io/docs/reference/project-config

## Development

To run in live development mode, run `wails dev` in the project directory. This will run a Vite development
server that will provide hot reload of your frontend changes. The app will also reload backend changes with live rebuild.

## Building

To build a production binary, run `wails build`. The output is generated in `build/bin/`.

Windows installer (NSIS): `wails build --nsis`

## Documentation

Complete project documentation is available in the `docs/` directory:

- **[INSTALLATION.md](docs/INSTALLATION.md)** — Platform-specific setup and troubleshooting
- **[USAGE.md](docs/USAGE.md)** — User guide and feature documentation  
- **[CHANGELOG.md](docs/CHANGELOG.md)** — Version history and release notes
- **[CONTRIBUTING.md](docs/CONTRIBUTING.md)** — How to contribute to the project
- **[MAINTAINER.md](docs/MAINTAINER.md)** — Maintainer notes and development practices

## License

See LICENSE file in the project root.
