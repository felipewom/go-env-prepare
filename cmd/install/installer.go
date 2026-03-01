package install

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"slices"
	"strings"
)

type Installer interface {
	Install()
	IsAlreadyInstalled() bool
	Title() string
	Description() string
}

type LifecycleInstaller interface {
	Plan() string
	Check() bool
	Apply() error
}

var installers = []Installer{
	&HomebrewInstaller{},
	&Iterm2Installer{},
	&ZshInstaller{},
	&VscodeInstaller{},
	&GitInstaller{},
	&GoInstaller{},
	&NodeJSInstaller{},
	&DotnetInstaller{},
	&PythonInstaller{},
	&DockerInstaller{},
}

func GetTitles() []string {
	var titles []string
	for _, i := range installers {
		if i.Title() == "" {
			continue
		}
		titles = append(titles, i.Title())
	}
	return titles
}

func GetInstallerByTitle(title string) Installer {
	for _, i := range installers {
		if i.Title() == title {
			return i
		}
	}
	return nil
}

func GetDescriptionByTitle(title string) string {
	installer := GetInstallerByTitle(title)
	if installer == nil {
		return ""
	}

	return installer.Description()
}

func GetAlreadyInstalledTools() []string {
	var alreadyInstalledTools []string
	for _, i := range installers {
		if !IsInstalled(i) {
			continue
		}
		alreadyInstalledTools = append(alreadyInstalledTools, i.Title())
	}
	return alreadyInstalledTools
}

func IsInstalled(i Installer) bool {
	if lifecycleInstaller, ok := i.(LifecycleInstaller); ok {
		return lifecycleInstaller.Check()
	}
	return i.IsAlreadyInstalled()
}

func RunInstaller(i Installer) {
	if lifecycleInstaller, ok := i.(LifecycleInstaller); ok {
		if lifecycleInstaller.Check() {
			return
		}
		if err := lifecycleInstaller.Apply(); err != nil {
			fmt.Printf("\r❌ Error installing %s: %v\n", i.Title(), err)
		}
		return
	}
	i.Install()
}

func FinishAllInstallations() {
	fmt.Printf("\r✅ Installation steps finished.\n")

	// Reload the shell configuration
	err := reloadShellConfiguration()
	if err != nil {
		fmt.Printf("\r⚠️ Install completed, but shell configuration was not reloaded automatically: %v\n", err)
		fmt.Printf("\r⚠️ Please run `source ~/.zshrc` (or restart your shell) to load changes.\n")
		return
	}

	fmt.Printf("\r✅ Shell configuration reloaded. Some tools may still require a new terminal session.\n")
}

func reloadShellConfiguration() error {
	// Get the user's home directory
	usr, err := user.Current()
	if err != nil {
		fmt.Printf("\r❌ Error getting user's home directory: %v\n", err)
		return err
	}

	shell := currentShellPath()
	if shell == "" {
		return errors.New("unable to determine active shell")
	}

	configFile, err := shellConfigPath(usr.HomeDir, shell)
	if err != nil {
		fmt.Println("Warning: Unsupported shell detected. Please update your shell configuration manually.")
		return err
	}

	// Execute the command to reload the shell configuration
	cmd := exec.Command(shell, "-c", fmt.Sprintf("source %s", configFile))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func shellConfigPath(homeDir, shellPath string) (string, error) {
	switch filepath.Base(shellPath) {
	case "zsh":
		return filepath.Join(homeDir, ".zshrc"), nil
	case "bash":
		return filepath.Join(homeDir, ".bashrc"), nil
	default:
		return "", errors.New("unsupported shell")
	}
}

func currentShellPath() string {
	if shell := os.Getenv("SHELL"); shell != "" {
		return shell
	}
	for _, candidate := range []string{"zsh", "bash"} {
		path, err := exec.LookPath(candidate)
		if err == nil {
			return path
		}
	}
	return ""
}

func shellListedInEtcShells(shellPath string) bool {
	content, err := os.ReadFile("/etc/shells")
	if err != nil {
		return false
	}
	allowed := parseShellList(string(content))
	return slices.Contains(allowed, shellPath)
}

func parseShellList(content string) []string {
	lines := strings.Split(content, "\n")
	allowed := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		allowed = append(allowed, line)
	}
	return allowed
}
