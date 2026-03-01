# Roadmap

This document tracks planned features and the release cadence for `go-env-prepare`.

## Release Cadence

| Cadence | Description |
|---------|-------------|
| **Patch** (x.x.N) | Bug fixes, doc updates — as needed |
| **Minor** (x.N.0) | New tool installers, UX improvements — monthly or on demand |
| **Major** (N.0.0) | Breaking interface changes — rare, announced in advance |

Releases are automated via `.github/workflows/release.yml`. Push a tag `vX.Y.Z` to trigger.

---

## Upcoming (next minor)

- [ ] `--non-interactive` / `--all` flags for scripted use
- [ ] `--profile` flag to load a predefined tool set (backend-go, frontend, etc.)
- [ ] Rust toolchain installer (`rustup`)
- [ ] Python version manager installer (`pyenv`)
- [ ] `prepare update` subcommand to update all installed tools

## Backlog

- [ ] Windows / WSL2 support
- [ ] Shell completion (`prepare completion zsh`)
- [ ] `prepare list` — show all tools and their install status
- [ ] Verbose mode (`--verbose`) for debugging install steps
- [ ] Plugin system for community-contributed installers

## Governance

- Roadmap is reviewed and updated with each minor release.
- Items move to "Upcoming" when there is an open PR or active development.
- Community requests are tracked via GitHub Issues with the `enhancement` label.

## Changelog

See [CHANGELOG.md](../CHANGELOG.md) for a history of releases.
