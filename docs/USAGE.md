Usage

- Launch the app from the Start menu / Applications folder or run the binary in `build/bin`.
- Profiles are stored under the app settings — use the UI to manage profiles.
- Session persistence: if an in-progress session is found (and not older than 24 hours) the app prompts to resume. The persistence file is `state.json` stored in the OS cache directory for the app.

Shortcuts

- `M` toggles mini timer (if implemented in frontend).

Troubleshooting

- If the app shows a previous session after reinstall, ensure you uninstalled the previous version (the uninstaller removes `%LOCALAPPDATA%\FocusPlay`).
- Editor warning about `wails.json` schema is cosmetic and has been removed from the repo to avoid VS Code warnings.
