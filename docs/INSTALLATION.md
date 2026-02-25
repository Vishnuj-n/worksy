Installation

Windows (recommended):

1. Build the installer:

   wails build --nsis

2. Run the generated installer: `build\bin\focusplay-amd64-installer.exe`.
3. To uninstall, use the Control Panel > Programs or run the uninstaller created by the installer. The uninstaller now removes user data at `%LOCALAPPDATA%\FocusPlay`.

macOS:

1. Build for macOS:

   wails build --platform darwin/universal

2. Open the produced .app or installer and install normally.

Linux:

1. Build for Linux:

   wails build --platform linux/amd64

2. Run the binary produced in `build/bin/`.
