# Release Guide

## Triggering a Release

Push a version tag to trigger the release workflow:

```bash
git tag v0.1.0
git push origin v0.1.0
```

The workflow (`.github/workflows/release.yml`) will:
1. Run `go test -race ./...`
2. Build binaries for darwin/linux × amd64/arm64
3. Generate `checksums.txt`
4. Publish a GitHub Release with all artifacts

## Homebrew Formula Update (Opt-in)

> ⚠️ **Formula update is disabled by default.**
>
> The current `formula.rb` install block uses a source tarball approach that
> is not yet validated end-to-end. Enabling it before the formula is correct
> will push a broken formula to users.

### How to enable

Once `formula.rb` has a working `install` block (e.g., using `bin.install`
against a pre-built binary tarball), set the following **repository variable**
in GitHub → Settings → Variables → Actions:

| Variable | Value |
|----------|-------|
| `ENABLE_FORMULA_UPDATE` | `true` |

With this variable set, every release tag push will:
1. Download the source tarball and compute its `sha256`.
2. Run `scripts/update-formula.sh <tag> <sha256>` to patch `formula.rb`.
3. Commit and push the updated formula back to `main`.

If the push fails for any reason other than "no changes", the workflow **fails
loudly** — it will not silently swallow the error.

### Manual update (current recommended path)

Until the formula install block is correct, update the formula manually after
each release:

```bash
VERSION=v0.1.0
curl -fsSL "https://github.com/felipewom/go-env-prepare/archive/refs/tags/${VERSION}.tar.gz" \
  -o release.tar.gz
SHA256=$(sha256sum release.tar.gz | awk '{print $1}')
./scripts/update-formula.sh "${VERSION}" "${SHA256}"
# review formula.rb, then commit and push
git add formula.rb
git commit -m "chore(release): update formula to ${VERSION}"
git push origin main
```

## Versioning

This project follows [Semantic Versioning](https://semver.org/):

- **Patch** (`x.x.N`): bug fixes, doc updates
- **Minor** (`x.N.0`): new installers, non-breaking features
- **Major** (`N.0.0`): breaking changes (rare, announced in advance)
