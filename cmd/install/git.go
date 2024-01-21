package install

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

type GitInstaller struct{}

func (g *GitInstaller) Install() {
	fmt.Println("Checking Git installation...")

	if g.IsAlreadyInstalled() {
		return
	}
	fmt.Println("Git is not installed. Installing...")

	// Display an animated hourglass loading indicator
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				fmt.Printf("\r⌛ Installing Git %s", loadingAnimation())
				time.Sleep(300 * time.Millisecond)
			}
		}
	}()

	// Install Git using Homebrew
	cmd := exec.Command("brew", "install", "git")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	// Stop the loading animation
	close(done)

	if err != nil {
		fmt.Printf("\r❌ Error installing Git: %v\n", err)
		return
	}

	fmt.Printf("\r✅ Git installed successfully.\n")
}

func (g *GitInstaller) Title() string {
	return "Git"
}

func (g *GitInstaller) Description() string {
	alreadyInstalled := ""
	if g.IsAlreadyInstalled() {
		alreadyInstalled = " (Already installed)"
	}
	return fmt.Sprintf("This option will install Git.%s", alreadyInstalled)
}

func (g *GitInstaller) IsAlreadyInstalled() bool {
	_, err := exec.LookPath("git")
	return err == nil
}
