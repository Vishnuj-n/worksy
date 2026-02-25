Contributing

Thanks for your interest in contributing to FocusPlay.

- Fork the repository and create a feature branch.
- Run the app locally and ensure changes build:

  wails build

- For UI changes, edit files under `frontend/` and rebuild the frontend before running the Go build.

- Tests: some services include Go unit tests under `internal/.../service_test.go`. Run `go test ./...` to run available tests.

- Keep commits small and focused. Include test or manual verification steps for your changes.

- For major changes, open an issue first to discuss the design.
