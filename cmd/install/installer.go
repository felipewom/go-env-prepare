package install

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

type Installer interface {
	Install()
	IsAlreadyInstalled() bool
	Title() string
	Description() string
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
		if !i.IsAlreadyInstalled() {
			continue
		}
		alreadyInstalledTools = append(alreadyInstalledTools, i.Title())
	}
	return alreadyInstalledTools
}

func FinishAllInstallations() {
	// Reload the shell configuration
	err := reloadShellConfiguration()
	if err != nil {
		fmt.Printf("\r❌ Error reloading shell configuration: %v\n", err)
		return
	}

	fmt.Printf("\r✅ All installations completed. Please restart your shell.\n")
}

func reloadShellConfiguration() error {
	// Get the user's shell
	shell := os.Getenv("SHELL")
	// Path to the shell configuration file (.zshrc or .bashrc)
	// Get the user's home directory
	usr, err := user.Current()
	if err != nil {
		fmt.Printf("\r❌ Error getting user's home directory: %v\n", err)
		return err
	}
	configFile := ""
	switch {
	case strings.Contains(shell, "zsh"):
		configFile = filepath.Join(usr.HomeDir, ".zshrc")
	case strings.Contains(shell, "bash"):
		configFile = filepath.Join(usr.HomeDir, ".bashrc")
	default:
		fmt.Println("Warning: Unsupported shell detected. Please update your shell configuration manually.")
		return errors.New("unsupported shell")
	}
	// Determine the command to reload the shell based on the shell type
	reloadCommand := "source"
	// Execute the command to reload the shell configuration
	cmd := exec.Command(os.Getenv("SHELL"), "-c", fmt.Sprintf("%s %s", reloadCommand, configFile))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
