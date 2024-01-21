package install

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

type VscodeInstaller struct{}

func (v *VscodeInstaller) Install() {
	fmt.Println("Checking Visual Studio Code installation...")

	if v.IsAlreadyInstalled() {
		return
	}
	fmt.Println("Visual Studio Code is not installed. Installing...")

	// Display an animated hourglass loading indicator
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				fmt.Printf("\r⌛ Installing Visual Studio Code %s", loadingAnimation())
				time.Sleep(300 * time.Millisecond)
			}
		}
	}()

	// Install Visual Studio Code using Homebrew
	cmd := exec.Command("brew", "install", "visual-studio-code")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	// Stop the loading animation
	close(done)

	if err != nil {
		fmt.Printf("\r❌ Error installing Visual Studio Code: %v\n", err)
		return
	}

	fmt.Printf("\r✅ Visual Studio Code installed successfully.\n")

	// Perform post-installation configurations
	v.postInstall()
}

func (v *VscodeInstaller) Title() string {
	return "Visual Studio Code"
}

func (v *VscodeInstaller) Description() string {
	alreadyInstalled := ""
	if v.IsAlreadyInstalled() {
		alreadyInstalled = " (Already installed)"
	}
	return fmt.Sprintf("This option will install Visual Studio Code.%s", alreadyInstalled)
}

// IsVscodeInstalled checks if Visual Studio Code is already installed
func (v *VscodeInstaller) IsAlreadyInstalled() bool {
	_, err := exec.LookPath("code")
	return err == nil
}

func (v *VscodeInstaller) postInstall() {
	// Set Visual Studio Code as the default Git editor
	v.setGitEditor()
}

func (v *VscodeInstaller) setGitEditor() {
	fmt.Println("Setting Visual Studio Code as the default Git editor...")

	// Configure Git to use Visual Studio Code as the editor
	cmd := exec.Command("git", "config", "--global", "core.editor", "code --wait")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		fmt.Printf("\r❌ Error setting Git editor: %v\n", err)
		return
	}

	fmt.Printf("\r✅ Visual Studio Code set as the default Git editor.\n")
}
