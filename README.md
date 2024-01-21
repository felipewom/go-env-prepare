# Go Env Prepare

![MIT License](https://img.shields.io/badge/license-MIT-green)
[![Buy Me a Coffee](https://img.shields.io/badge/buy%20me%20a%20coffee-donate-orange.svg)](https://buymeacoffee.com/felipewom)

This is a CLI library that prepare your environment for different stacks, such as Node.js, Go, React, and .NET
Your CLI Project is a command-line tool to help you set up development environments for various stacks.

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
