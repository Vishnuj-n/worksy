Changelog

## 2026-02-25 - Unreleased

- Renamed Wails output filename to `focusplay`.
- Removed external JSON schema line from `wails.json` to silence editor warnings.
- Uninstaller now removes user data directory `%LOCALAPPDATA%\FocusPlay` to avoid stale session restore after reinstall.
- Rebuilt installer with `wails build --nsis`.
