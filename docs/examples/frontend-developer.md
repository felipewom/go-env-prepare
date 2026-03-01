# Example: Frontend Developer

A profile for engineers building web applications with Node.js/npm on macOS.

## Recommended Tools

| Tool     | Reason |
|----------|--------|
| Homebrew | Package manager — required for most other installs |
| Git      | Version control |
| NodeJS   | JavaScript runtime (includes npm) |
| VS Code  | Editor with excellent JS/TS support |
| Zsh      | Default macOS shell; works well with `nvm` and `oh-my-zsh` |
| iTerm2   | Enhanced terminal with better split-pane and search support |

## Setup Walkthrough

1. Run `prepare` and select the tools above.
2. Install `nvm` for managing multiple Node versions:
   ```bash
   curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash
   # then reload your shell
   nvm install --lts
   nvm use --lts
   ```
3. Verify Node:
   ```bash
   node --version
   npm --version
   ```
4. Install recommended VS Code extensions:
   - ESLint, Prettier, GitLens, TypeScript Hero

## Typical Project Workflow

```bash
# Bootstrap a React/Next.js project
npx create-next-app@latest my-app
cd my-app

# Install dependencies
npm install

# Development server
npm run dev

# Run tests
npm test

# Build for production
npm run build
```

## Additional Recommended Tools (manual install)

```bash
brew install pnpm     # faster npm alternative
brew install volta    # alternative Node version manager
npm install -g typescript ts-node   # TypeScript tooling
```
