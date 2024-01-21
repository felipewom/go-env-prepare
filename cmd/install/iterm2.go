package install

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

type Iterm2Installer struct{}

func (i *Iterm2Installer) Install() {
	fmt.Println("Checking iTerm2 installation...")

	if i.IsAlreadyInstalled() {
		fmt.Println("iTerm2 is already installed.")
		return
	}

	fmt.Println("iTerm2 is not installed. Installing...")

	// Display an animated hourglass loading indicator
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				fmt.Printf("\r⌛ Installing iTerm2 %s", loadingAnimation())
				time.Sleep(300 * time.Millisecond)
			}
		}
	}()

	// Install iTerm2 using Homebrew
	cmd := exec.Command("brew", "install", "iterm2")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	// Stop the loading animation
	close(done)

	if err != nil {
		fmt.Printf("\r❌ Error installing iTerm2: %v\n", err)
		return
	}

	fmt.Printf("\r✅ iTerm2 installed successfully.\n")

	// Set iTerm2 as the default terminal
	i.setDefaultTerminal()
}

func (i *Iterm2Installer) Title() string {
	return "iTerm2"
}

func (i *Iterm2Installer) Description() string {
	alreadyInstalled := ""
	if i.IsAlreadyInstalled() {
		alreadyInstalled = " (Already installed)"
	}
	return fmt.Sprintf("This option will install iTerm2 and set it as the default terminal emulator.%s", alreadyInstalled)
}

func (i *Iterm2Installer) IsAlreadyInstalled() bool {
	_, err := os.Stat("/Applications/iTerm.app")
	_, err2 := exec.LookPath("iterm")
	return err == nil || err2 == nil
}

func (i *Iterm2Installer) setDefaultTerminal() {
	fmt.Println("Setting iTerm2 as the default terminal...")

	// Get the user's home directory
	usr, err := user.Current()
	if err != nil {
		fmt.Printf("\r❌ Error getting user's home directory: %v\n", err)
		return
	}

	// Detect the user's shell
	shell := os.Getenv("SHELL")
	var configFile string

	switch {
	case strings.Contains(shell, "zsh"):
		configFile = filepath.Join(usr.HomeDir, ".zshrc")
	case strings.Contains(shell, "bash"):
		configFile = filepath.Join(usr.HomeDir, ".bashrc")
	default:
		fmt.Println("Warning: Unsupported shell detected. Please update your shell configuration manually.")
		return
	}

	// Update the shell configuration file to set iTerm2 as the default terminal
	cmd := exec.Command("bash", "-c", fmt.Sprintf(`echo 'export TERMINAL="%s"' >> %s`, "iterm", configFile))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()

	if err != nil {
		fmt.Printf("\r❌ Error setting iTerm2 as the default terminal: %v\n", err)
		return
	}

	fmt.Printf("\r✅ iTerm2 set as the default terminal. Please restart your shell.\n")
}
