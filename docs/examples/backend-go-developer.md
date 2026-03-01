# Example: Backend Go Developer

A profile for engineers building Go services, CLIs, or APIs on macOS.

## Recommended Tools

| Tool     | Reason |
|----------|--------|
| Homebrew | Package manager — required for most other installs |
| Git      | Version control |
| Go       | Primary language runtime |
| Docker   | Local container runtime for services and integration tests |
| VS Code  | Editor with Go extension support |
| Zsh      | Default macOS shell; required for `oh-my-zsh` ecosystem |

## Setup Walkthrough

1. Run `prepare` and select the tools above.
2. After completion, install the official Go extension in VS Code:
   - Open VS Code → `Cmd+Shift+X` → search "Go" → Install
3. Verify your Go environment:
   ```bash
   go version
   go env GOPATH
   ```
4. Verify Docker:
   ```bash
   docker run --rm hello-world
   ```

## Typical Project Workflow

```bash
# Bootstrap a new project
mkdir my-service && cd my-service
go mod init github.com/you/my-service

# Run tests
go test -race ./...

# Build and run in Docker
docker build -t my-service .
docker run --rm -p 8080:8080 my-service
```

## Additional Recommended Tools (manual install)

```bash
brew install golangci-lint   # linter
brew install goreleaser      # release tooling
brew install air             # hot reload for development
```
