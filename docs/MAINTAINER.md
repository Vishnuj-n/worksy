MAINTAINER NOTES

Repository: FocusPlay (formerly Worksy output filename fixed)
Date: 2026-02-25
Maintainer: Vishnu J Narayanan

Summary of recent changes

- Renamed the Wails output filename to `focusplay` in `wails.json` (was `Worksy`).
- Removed the external JSON schema line from `wails.json` to avoid VS Code editor warnings.
- Added uninstaller cleanup to the NSIS installer (`build/windows/installer/project.nsi`) to remove the user data directory (`%LOCALAPPDATA%\FocusPlay`) during uninstall. This prevents stale session state from persisting across re-installs.
- Confirmed persistence lives under the OS cache directory via `internal/infra/storage/json_store.go` which uses `os.UserCacheDir()`.
- Built a new installer using `wails build --nsis` which produced `build/bin/focusplay.exe` and the NSIS installer.

Notes for maintainers

- Where state is stored: `internal/services/persistence/service.go` writes `state.json` to the data dir returned by `storage.DataDir()`.
- Session expiry: the persistence service treats state older than 24 hours as stale and will ignore it on load.

Recommended next steps

- Consider adding an explicit "Reset data" option in the app settings to allow users to clear persisted state without uninstalling.
- Add automated installer tests to ensure uninstall removes user data on Windows.

Build & test commands

- Build the app and NSIS installer (Windows):

  wails build --nsis

- Run the produced installer:

  Start-Process "C:\\Users\\vishn\\PROJECT\\POMODORO\\build\\bin\\focusplay-amd64-installer.exe"

Contact

For issues ask the original author: Vishnu J Narayanan <144352066+Vishnuj-n@users.noreply.github.com>
