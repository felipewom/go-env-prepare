# Go Env Prepare

![MIT License](https://img.shields.io/badge/license-MIT-green)
[![Buy Me a Coffee](https://img.shields.io/badge/buy%20me%20a%20coffee-donate-orange.svg)](https://buymeacoffee.com/felipewom)

This is a CLI library that prepares your environment for different stacks, such as Node.js, Go, React, and .NET.
It now supports both:
- Interactive installer flow (`prepare`)
- Declarative dynamic engine (`prepare plan|run|lint|lock`)

## Getting Started

```markdown

go-env-prepare/
|-- cmd/
|   |-- root.go
|   |-- commands.go
|   |-- install/
|       |-- go.go
|       |-- homebrew.go
|       |-- nodejs.go
|       |-- dotnet.go
|       |-- docker.go
|-- main.go
|-- go.mod
|-- go.sum
|-- Makefile
|-- README.md

```

### Prerequisites

- Go installed
- Docker installed
- Node.js installed

## Installation

To install the CLI, run the following commands:

```bash
make install
```

## Usage

To start a new Go project with Go Modules, run:

```bash
make run
```

Select the desired development tools when prompted.

### Dynamic, Declarative Execution

Use a manifest (`prepare.yaml` or `prepare.json`) to define profiles and tools. If no manifest is provided, builtin profiles are used.

Example `prepare.yaml`:

```yaml
apiVersion: v1
profile: fullstack
profiles:
  my-stack:
    extends:
      - backend
    tools:
      - nodejs
```

Commands:

```bash
# Build dependency-aware plan
prepare plan --profile frontend

# Show exact execution without changing machine
prepare run --profile backend --dry-run

# Load user manifest and execute
prepare run --file prepare.yaml --profile my-stack

# Validate profile syntax and semantics
prepare lint --file prepare.yaml

# Structured output for automation
prepare run --profile fullstack --dry-run --json
```

Architecture summary:
- Manifest loader/parser: resolves builtin + user profiles with inheritance.
- Planner: expands dependencies and generates a topological execution order.
- Executor: idempotent checks (`already_installed` skip), dry-run mode, and checkpoint resume (`--resume --state`).
- Lock strategy: generates `prepare.lock.json` with pinned source/version metadata from the resolved plan.

Migration notes from static installer flow:
- `prepare` (no subcommand) keeps the existing interactive experience.
- New workflows are additive under `prepare plan|run|lint|lock`.
- Existing installer files under `cmd/install` are preserved for backward compatibility.

## Build

To build the CLI binary, run:

```bash
make build
```

The binary will be named `your-cli`.

## Cleaning Up

To clean the project, run:

```bash
make clean
```

This will remove the binary.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
