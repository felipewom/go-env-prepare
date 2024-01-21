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

type NodeJSInstaller struct{}

func (n *NodeJSInstaller) Install() {
	fmt.Println("Checking NVM installation...")

	if n.IsAlreadyInstalled() {
		return
	}

	fmt.Println("NVM is not installed. Installing...")

	// Display an animated hourglass loading indicator
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				fmt.Printf("\r⌛ Installing NVM %s", loadingAnimation())
				time.Sleep(300 * time.Millisecond)
			}
		}
	}()

	// Install NVM using Homebrew
	cmd := exec.Command("brew", "install", "nvm")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	// Stop the loading animation
	close(done)

	if err != nil {
		fmt.Printf("\r❌ Error installing NVM: %v\n", err)
		return
	}

	fmt.Printf("\r✅ NVM installed successfully.\n")

	// Install the latest LTS version of Node.js
	n.installLatestLTSNode()

	// Set NVM environment variables
	n.setNvmVariables()
}

func (n *NodeJSInstaller) Title() string {
	return "NVM (NodeJS LTS)"
}

func (n *NodeJSInstaller) Description() string {
	alreadyInstalled := ""
	if n.IsAlreadyInstalled() {
		alreadyInstalled = " (Already installed)"
	}
	return fmt.Sprintf("This option will install nvm and the latest LTS version of Node.js.%s", alreadyInstalled)
}

func (n *NodeJSInstaller) IsAlreadyInstalled() bool {
	_, exists := os.LookupEnv("NVM_DIR")
	return exists
}

func (n *NodeJSInstaller) installLatestLTSNode() {
	fmt.Println("Installing the latest LTS version of Node.js...")

	// Source NVM in the current shell
	cmd := exec.Command("bash", "-c", "source $(brew --prefix nvm)/nvm.sh && nvm install --lts")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		fmt.Printf("\r❌ Error installing Node.js: %v\n", err)
		return
	}

	fmt.Printf("\r✅ Node.js installed successfully.\n")
}

func (n *NodeJSInstaller) setNvmVariables() {
	fmt.Println("Setting NVM environment variables...")

	// Get the user's shell
	shell := os.Getenv("SHELL")

	// Get the user's home directory
	usr, err := user.Current()
	if err != nil {
		fmt.Printf("\r❌ Error getting user's home directory: %v\n", err)
		return
	}

	// Path to the shell configuration file (.zshrc or .bashrc)
	configFile := ""
	switch {
	case strings.Contains(shell, "zsh"):
		configFile = filepath.Join(usr.HomeDir, ".zshrc")
	case strings.Contains(shell, "bash"):
		configFile = filepath.Join(usr.HomeDir, ".bashrc")
	default:
		fmt.Println("Warning: Unsupported shell detected. Please update your shell configuration manually.")
		return
	}

	// NVM initialization script
	nvmScript := `

# NVM initialization script
export NVM_DIR="$HOME/.nvm"
[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"  # This loads nvm
[ -s "$NVM_DIR/bash_completion" ] && \. "$NVM_DIR/bash_completion"  # This loads nvm bash_completion
	`

	// Append the NVM initialization script to the shell configuration file
	f, err := os.OpenFile(configFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("\r❌ Error opening shell configuration file: %v\n", err)
		return
	}
	defer f.Close()

	_, err = f.WriteString(nvmScript)
	if err != nil {
		fmt.Printf("\r❌ Error writing to shell configuration file: %v\n", err)
		return
	}

	fmt.Printf("\r✅ NVM environment variables set successfully. Please restart your shell.\n")
}
