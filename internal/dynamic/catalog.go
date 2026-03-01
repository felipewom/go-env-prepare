package dynamic

func BuiltinProfiles() map[string]Profile {
	return map[string]Profile{
		"base": {
			Tools: []string{"homebrew", "git", "zsh", "vscode"},
		},
		"frontend": {
			Extends: []string{"base"},
			Tools:   []string{"nodejs"},
		},
		"backend": {
			Extends: []string{"base"},
			Tools:   []string{"go", "python", "docker"},
		},
		"data": {
			Extends: []string{"base"},
			Tools:   []string{"python", "docker"},
		},
		"ai": {
			Extends: []string{"data"},
			Tools:   []string{"dotnet"},
		},
		"fullstack": {
			Extends: []string{"frontend", "backend"},
			Tools:   []string{"dotnet", "iterm2"},
		},
	}
}

func BuiltinCatalog() map[string]ToolSpec {
	return map[string]ToolSpec{
		"homebrew": {
			ID:          "homebrew",
			Title:       "Homebrew",
			Description: "Package manager for macOS",
			Install: Command{
				Name:  "/bin/bash",
				Args:  []string{"-c", "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"},
				Shell: "bash",
			},
			Check:   Check{Binary: "brew"},
			Version: "latest",
			Source:  "homebrew/homebrew-core",
		},
		"iterm2": {
			ID:           "iterm2",
			Title:        "iTerm2",
			Description:  "Terminal emulator",
			Dependencies: []string{"homebrew"},
			Install:      Command{Name: "brew", Args: []string{"install", "--cask", "iterm2"}},
			Check:        Check{PathExists: "/Applications/iTerm.app"},
			Version:      "latest",
			Source:       "homebrew/cask",
		},
		"zsh": {
			ID:           "zsh",
			Title:        "Zsh",
			Description:  "Shell",
			Dependencies: []string{"homebrew"},
			Install:      Command{Name: "brew", Args: []string{"install", "zsh"}},
			Check:        Check{Binary: "zsh"},
			Version:      "latest",
			Source:       "homebrew/homebrew-core",
		},
		"vscode": {
			ID:           "vscode",
			Title:        "Visual Studio Code",
			Description:  "Code editor",
			Dependencies: []string{"homebrew"},
			Install:      Command{Name: "brew", Args: []string{"install", "--cask", "visual-studio-code"}},
			Check:        Check{Binary: "code"},
			Version:      "latest",
			Source:       "homebrew/cask",
		},
		"git": {
			ID:           "git",
			Title:        "Git",
			Description:  "Version control",
			Dependencies: []string{"homebrew"},
			Install:      Command{Name: "brew", Args: []string{"install", "git"}},
			Check:        Check{Binary: "git"},
			Version:      "latest",
			Source:       "homebrew/homebrew-core",
		},
		"go": {
			ID:           "go",
			Title:        "Go",
			Description:  "Go programming language",
			Dependencies: []string{"homebrew"},
			Install:      Command{Name: "brew", Args: []string{"install", "go"}},
			Check:        Check{Binary: "go"},
			Version:      "latest",
			Source:       "homebrew/homebrew-core",
		},
		"nodejs": {
			ID:           "nodejs",
			Title:        "Node.js",
			Description:  "Node.js LTS runtime",
			Dependencies: []string{"homebrew"},
			Install:      Command{Name: "brew", Args: []string{"install", "node"}},
			Check:        Check{Binary: "node"},
			Version:      "lts",
			Source:       "homebrew/homebrew-core",
		},
		"dotnet": {
			ID:           "dotnet",
			Title:        ".NET SDK",
			Description:  ".NET SDK",
			Dependencies: []string{"homebrew"},
			Install:      Command{Name: "brew", Args: []string{"install", "dotnet-sdk"}},
			Check:        Check{Binary: "dotnet"},
			Version:      "latest",
			Source:       "homebrew/homebrew-core",
		},
		"python": {
			ID:           "python",
			Title:        "Python",
			Description:  "Python runtime",
			Dependencies: []string{"homebrew"},
			Install:      Command{Name: "brew", Args: []string{"install", "python"}},
			Check:        Check{Binary: "python3"},
			Version:      "latest",
			Source:       "homebrew/homebrew-core",
		},
		"docker": {
			ID:           "docker",
			Title:        "Docker",
			Description:  "Container runtime",
			Dependencies: []string{"homebrew"},
			Install:      Command{Name: "brew", Args: []string{"install", "docker"}},
			Check:        Check{Binary: "docker"},
			Version:      "latest",
			Source:       "homebrew/homebrew-core",
		},
	}
}
