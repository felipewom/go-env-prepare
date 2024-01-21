package install

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

type GoInstaller struct{}

func (g *GoInstaller) Install() {
	fmt.Println("Checking Go installation...")

	if g.IsAlreadyInstalled() {
		return
	}

	fmt.Println("Go is not installed. Installing...")

	// Display an animated hourglass loading indicator
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				fmt.Printf("\r⌛ Installing Go %s", loadingAnimation())
				time.Sleep(300 * time.Millisecond)
			}
		}
	}()

	// Install Go using Homebrew
	cmd := exec.Command("brew", "install", "go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	// Stop the loading animation
	close(done)

	if err != nil {
		fmt.Printf("\r❌ Error installing Go: %v\n", err)
		return
	}

	fmt.Printf("\r✅ Go installed successfully.\n")
}

func (g *GoInstaller) Title() string {
	return "Go"
}

func (g *GoInstaller) IsAlreadyInstalled() bool {
	_, err := exec.LookPath("go")
	return err == nil
}

func (g *GoInstaller) Description() string {
	alreadyInstalled := ""
	if g.IsAlreadyInstalled() {
		alreadyInstalled = " (Already installed)"
	}
	return fmt.Sprintf("This option will install latest version of Go.%s", alreadyInstalled)
}
