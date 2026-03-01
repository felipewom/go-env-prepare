# Contributing to go-env-prepare

Thanks for your interest in contributing! This guide covers everything you need to get started.

## Prerequisites

- Go 1.21+
- macOS (primary target platform)
- `git`

## Development Setup

```bash
git clone https://github.com/felipewom/go-env-prepare.git
cd go-env-prepare
go mod tidy
go build -o prepare .
```

## Running Tests

```bash
go test ./...
go test -race ./...          # with race detector
go test -v -run TestDetection ./cmd/install/  # specific test
```

## Making Changes

1. **Fork** the repo and create a branch: `git checkout -b feat/my-feature`
2. Make your change with tests.
3. Run `go test ./...` — all tests must pass.
4. Commit using the format below.

## Commit Message Format

```
type(scope): short description

Longer body if needed.
```

**Types:** `feat`, `fix`, `chore`, `docs`, `ci`, `refactor`, `test`  
**Examples:**
```
feat(install): add Rust toolchain installer
fix(homebrew): handle Rosetta 2 detection on M1
docs(oss): add troubleshooting section
ci: add race-detector check to CI matrix
```

## Code Style

- Keep installers **idempotent** — safe to run multiple times.
- Prefer **explicit error messaging** over silent failures.
- Commands should be **non-interactive** where possible.
- Follow standard Go formatting (`gofmt`).

## Adding a New Tool Installer

1. Create `cmd/install/<tool>.go` implementing the `Installer` interface:
   ```go
   type MyToolInstaller struct{}
   func (m *MyToolInstaller) Title() string       { return "MyTool" }
   func (m *MyToolInstaller) Description() string { return "..." }
   func (m *MyToolInstaller) IsAlreadyInstalled() bool { ... }
   func (m *MyToolInstaller) Install()            { ... }
   ```
2. Register it in `cmd/install/installer.go` in the `installers` slice.
3. Add tests in `cmd/install/installer_test.go` (see `TestDetection_WithFake*` examples).

## Pull Request Process

- Open a PR against `main`.
- Fill in the PR template.
- One approving review required before merge.

## Reporting Issues

Use the issue templates in `.github/ISSUE_TEMPLATE/`.

## Code of Conduct

This project follows the [Contributor Covenant](CODE_OF_CONDUCT.md).
