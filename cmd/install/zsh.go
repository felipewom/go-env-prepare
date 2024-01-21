package install

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"time"
)

type ZshInstaller struct{}

func (z *ZshInstaller) Install() {
	fmt.Println("Checking Zsh installation...")

	if z.IsAlreadyInstalled() {
		return
	}
	fmt.Println("Zsh is not installed. Installing...")

	// Display an animated hourglass loading indicator
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				fmt.Printf("\r⌛ Installing Zsh %s", loadingAnimation())
				time.Sleep(300 * time.Millisecond)
			}
		}
	}()

	// Install Zsh using Homebrew
	cmd := exec.Command("brew", "install", "zsh")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	// Stop the loading animation
	close(done)

	if err != nil {
		fmt.Printf("\r❌ Error installing Zsh: %v\n", err)
		return
	}

	fmt.Printf("\r✅ Zsh installed successfully.\n")

	// Set Zsh as the default shell
	z.setDefaultShell()

	// Install Oh My Zsh
	z.installOhMyZsh()

	// Perform post-installation tasks for Zsh
	z.postInstall()
}

func (z *ZshInstaller) Title() string {
	return "Zsh"
}

func (z *ZshInstaller) Description() string {
	alreadyInstalled := ""
	if z.IsAlreadyInstalled() {
		alreadyInstalled = " (Already installed)"
	}
	return fmt.Sprintf("This option will install Zsh and set it as the default shell.%s", alreadyInstalled)
}

func (z *ZshInstaller) IsAlreadyInstalled() bool {
	_, err := exec.LookPath("zsh")
	return err == nil
}

// IsZshInstalled checks if Zsh is already installed
func IsZshInstalled() bool {
	_, err := exec.LookPath("zsh")
	return err == nil
}

func (z *ZshInstaller) setDefaultShell() {
	fmt.Println("Setting Zsh as the default shell...")

	// Get the user's home directory
	usr, err := user.Current()
	if err != nil {
		fmt.Printf("\r❌ Error getting user's home directory: %v\n", err)
		return
	}

	// Update the shell to Zsh
	cmd := exec.Command("chsh", "-s", "/usr/local/bin/zsh", usr.Username)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()

	if err != nil {
		fmt.Printf("\r❌ Error setting Zsh as the default shell: %v\n", err)
		return
	}

	fmt.Printf("\r✅ Zsh set as the default shell. Please restart your shell.\n")
}

func (z *ZshInstaller) installOhMyZsh() {
	fmt.Println("Installing Oh My Zsh...")

	// Install Oh My Zsh using the official script
	cmd := exec.Command("sh", "-c", "$(curl -fsSL https://raw.github.com/ohmyzsh/ohmyzsh/master/tools/install.sh)")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		fmt.Printf("\r❌ Error installing Oh My Zsh: %v\n", err)
		return
	}

	fmt.Printf("\r✅ Oh My Zsh installed successfully.\n")
}

func (z *ZshInstaller) postInstall() {
	fmt.Println("Performing post-installation tasks for Zsh...")

	// Install Zsh plugins (e.g., zsh-autosuggestions, zsh-syntax-highlighting)
	z.installZshPlugins()

	// Update the .zshrc file with additional configurations
	z.updateZshrc()
}

func (z *ZshInstaller) installZshPlugins() {
	fmt.Println("Installing Zsh plugins...")

	// Install zsh-autosuggestions
	cmd := exec.Command("git", "clone", "https://github.com/zsh-users/zsh-autosuggestions.git", "$ZSH_CUSTOM/plugins/zsh-autosuggestions")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("\r❌ Error installing zsh-autosuggestions: %v\n", err)
		return
	}

	// Install zsh-syntax-highlighting
	cmd = exec.Command("git", "clone", "https://github.com/zsh-users/zsh-syntax-highlighting.git", "$ZSH_CUSTOM/plugins/zsh-syntax-highlighting")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Printf("\r❌ Error installing zsh-syntax-highlighting: %v\n", err)
		return
	}

	fmt.Printf("\r✅ Zsh plugins installed successfully.\n")
}

func (z *ZshInstaller) updateZshrc() {
	fmt.Println("Updating .zshrc file...")

	// Get the user's home directory
	usr, err := user.Current()
	if err != nil {
		fmt.Printf("\r❌ Error getting user's home directory: %v\n", err)
		return
	}

	// Path to the .zshrc file
	zshrcPath := filepath.Join(usr.HomeDir, ".zshrc")

	// Append plugin configurations to the .zshrc file
	pluginsConfig := `
	# Zsh plugins
	plugins=(
		git
		docker
		nvm
		go
		zsh-autosuggestions
		zsh-syntax-highlighting
	)

	alias zc="code ~/.zshrc"
	alias rr="source ~/.zshrc && clear && echo 'Zsh configuration reloaded.'

	`
	f, err := os.OpenFile(zshrcPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("\r❌ Error opening .zshrc file: %v\n", err)
		return
	}
	defer f.Close()

	_, err = f.WriteString(pluginsConfig)
	if err != nil {
		fmt.Printf("\r❌ Error writing to .zshrc file: %v\n", err)
		return
	}

	fmt.Printf("\r✅ .zshrc file updated successfully.\n")
}
