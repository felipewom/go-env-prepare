# Changelog

All notable changes to `go-env-prepare` are documented here.

Format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).  
Versioning follows [Semantic Versioning](https://semver.org/).

---

## [Unreleased]

### Added
- Unit and integration tests for installer detection logic (ticket-028, ticket-029)
- CI workflow: lint / test (race) / build matrix on ubuntu + macOS (ticket-030)
- Release workflow: multi-arch binaries + sha256 checksums (ticket-031)
- Automated Homebrew formula update script `scripts/update-formula.sh` (ticket-032)
- `CONTRIBUTING.md`, `CODE_OF_CONDUCT.md`, issue/PR templates (ticket-033)
- Architecture documentation `docs/architecture.md` (ticket-034)
- Troubleshooting handbook `docs/troubleshooting.md` (ticket-035)
- Examples catalog: backend-go-developer, frontend-developer (ticket-036)
- `SECURITY.md` and Dependabot configuration (ticket-037)
- Roadmap and changelog cadence doc (ticket-040)

---

## [0.0.1] — 2024-01-01

### Added
- Initial release
- Interactive multi-select prompt for tool installation
- Installers: Homebrew, iTerm2, Zsh, VS Code, Git, Go, NodeJS, .NET, Python, Docker
- `prepare version` subcommand
- Homebrew formula (`formula.rb`)

[Unreleased]: https://github.com/felipewom/go-env-prepare/compare/0.0.1...HEAD
[0.0.1]: https://github.com/felipewom/go-env-prepare/releases/tag/0.0.1
