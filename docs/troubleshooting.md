# Troubleshooting

Common issues and their fixes when using `go-env-prepare`.

---

## `prepare: command not found`

**Cause:** The binary is not on your PATH.

**Fix:**
```bash
# If installed via Homebrew
brew link prepare

# If built from source
go build -o prepare .
sudo cp prepare /usr/local/bin/prepare  # or add build dir to PATH
```

---

## Homebrew installation fails (M1/M2 Mac)

**Cause:** Architecture detection for Rosetta 2 may produce unexpected results.

**Fix:**
```bash
# Install Homebrew manually for Apple Silicon
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
# Then add to PATH
echo 'eval "$(/opt/homebrew/bin/brew shellenv)"' >> ~/.zprofile
eval "$(/opt/homebrew/bin/brew shellenv)"
```

---

## `Error reloading shell configuration: unsupported shell`

**Cause:** Your shell is not `zsh` or `bash` (e.g., `fish`, `nu`).

**Fix:** `go-env-prepare` can still install tools; only the auto-reload step is skipped. Manually source your config:
```bash
source ~/.config/fish/config.fish  # fish
```

---

## Tool shows as "Already installed" but is broken

**Cause:** The binary exists on PATH but is corrupted or misconfigured.

**Fix:** Remove the broken binary and re-run `prepare`:
```bash
brew uninstall <tool>
brew install <tool>
# or
which <tool>   # find the path
rm $(which <tool>)
# then run prepare again
```

---

## Prompt shows no options / blank screen

**Cause:** Terminal does not support ANSI escape codes or pseudo-TTY.

**Fix:** Run `prepare` in a supported terminal (Terminal.app, iTerm2, Warp, VS Code integrated terminal).

---

## `go test ./...` fails locally

**Cause:** A test modifies `PATH` and leaks state (rare).

**Fix:**
```bash
go test -count=1 -p=1 ./...  # disable parallel + disable test caching
```

All `TestDetection_*` tests use `t.Cleanup` to restore `PATH`, so they should be safe. If you see flakiness, open a bug report.

---

## Release artifacts have wrong version in binary

**Cause:** `ldflags` not passed during build.

**Fix:** Build with:
```bash
go build -ldflags "-s -w -X main.version=v0.1.0" -o prepare .
```

---

## Formula sha256 mismatch after release

**Cause:** `formula.rb` was not updated after the release.

**Fix:**
```bash
TARBALL_URL="https://github.com/felipewom/go-env-prepare/archive/refs/tags/v0.1.0.tar.gz"
curl -fsSL "${TARBALL_URL}" -o release.tar.gz
SHA256=$(sha256sum release.tar.gz | awk '{print $1}')
./scripts/update-formula.sh v0.1.0 "${SHA256}"
```
