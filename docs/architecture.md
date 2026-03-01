# Architecture

This document describes the internals of `go-env-prepare` for contributors and users who want to extend it.

## Overview

```
main.go
  └─ cmd/Execute()
       └─ NewCommands()
            ├─ newInstallCmd()   ← root command (interactive prompt)
            └─ newVersionCmd()   ← `prepare version`
```

`go-env-prepare` is a **CLI tool** built with [Cobra](https://github.com/spf13/cobra) and [Survey](https://github.com/AlecAivazis/survey). Its single responsibility is to detect which development tools are already installed on macOS and install any that are missing.

---

## Key Concepts

### Installer Interface (`cmd/install/installer.go`)

Every supported tool implements the `Installer` interface:

```go
type Installer interface {
    Install()             // idempotent install logic
    IsAlreadyInstalled()  // returns true if tool binary is already on PATH
    Title() string        // display name shown in the prompt
    Description() string  // help text shown alongside the option
}
```

All installers are registered in the `installers` slice in `installer.go`. The order of registration determines the order shown in the interactive prompt.

### Detection Strategy

Each installer uses `exec.LookPath("<binary>")` to determine whether a tool is already installed. This means:
- Detection is PATH-based, not version-based.
- Installers are safe to call repeatedly (idempotent).

### Install Strategy

Most tools are installed via **Homebrew** (`brew install <package>`). Homebrew itself is installed via its official shell script. After all installations, the shell configuration file (`.zshrc` or `.bashrc`) is sourced automatically.

### Interactive Prompt

The root command presents a `survey.MultiSelect` prompt listing all registered installers. Tools already detected on PATH are pre-selected as defaults.

---

## Extension Points

### Adding a New Tool

1. Create `cmd/install/<tool>.go` implementing `Installer`.
2. Add it to the `installers` slice in `cmd/install/installer.go`.
3. Add detection and install tests in `cmd/install/installer_test.go`.

### Changing Install Behavior

Override the `Install()` method. Keep it idempotent — check `IsAlreadyInstalled()` at the top and return early if already present.

### Non-Homebrew Installers

Some tools (e.g., Homebrew itself) use alternative install strategies. See `cmd/install/homebrew.go` for an example of a curl-based install.

---

## Module Structure

```
go-env-prepare/
├── main.go                    # entry point
├── cmd/
│   ├── root.go                # Execute() entry
│   ├── commands.go            # Cobra command wiring + prompt loop
│   └── install/
│       ├── installer.go       # Installer interface + registry
│       ├── installer_test.go  # unit + integration tests
│       ├── animation.go       # loading animation helper
│       ├── homebrew.go
│       ├── git.go
│       ├── go.go
│       ├── nodejs.go
│       ├── docker.go
│       ├── dotnet.go
│       ├── iterm2.go
│       ├── vscode.go
│       ├── zsh.go
│       └── pyhton.go
├── scripts/
│   └── update-formula.sh      # Homebrew formula updater
├── .github/
│   ├── workflows/
│   │   ├── ci.yml             # lint / test / build
│   │   └── release.yml        # multi-arch release + formula update
│   ├── dependabot.yml
│   ├── ISSUE_TEMPLATE/
│   └── PULL_REQUEST_TEMPLATE.md
├── docs/
│   ├── architecture.md        # this file
│   ├── troubleshooting.md
│   ├── roadmap.md
│   └── examples/
│       ├── backend-go-developer.md
│       └── frontend-developer.md
└── formula.rb                 # Homebrew formula
```
