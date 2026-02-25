# README

## About

This is the official Wails Vanilla template.

You can configure the project by editing `wails.json`. More information about the project settings can be found
here: https://wails.io/docs/reference/project-config

## Live Development

To run in live development mode, run `wails dev` in the project directory. This will run a Vite development
server that will provide very fast hot reload of your frontend changes. If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect
to this in your browser, and you can call your Go code from devtools.

## Building

To build a redistributable, production mode package, use `wails build`.
## Documentation

Project documentation and maintainer notes are in the `docs/` directory. Key docs:

- `docs/MAINTAINER.md`: recent maintainer actions and notes.
- `docs/INSTALLATION.md`: platform-specific installation instructions.
- `docs/USAGE.md`: basic usage and troubleshooting.
- `docs/CONTRIBUTING.md`: how to contribute.
- `docs/CHANGELOG.md`: change history.

Please read `docs/MAINTAINER.md` for a summary of recent fixes (app output filename, uninstaller cleanup, and schema removal).
