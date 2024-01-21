package install

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type HomebrewInstaller struct{}

func (h *HomebrewInstaller) Install() {
	fmt.Println("Checking Homebrew installation...")

	if homebrewInstalled() {
		return
	}

	fmt.Println("Homebrew is not installed. Installing...")

	// Display an animated hourglass loading indicator
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				fmt.Printf("\r⌛ Installing Homebrew %s", loadingAnimation())
				time.Sleep(300 * time.Millisecond)
			}
		}
	}()

	var cmd *exec.Cmd

	// Determine the architecture for Homebrew installation
	arch := runtime.GOARCH
	if runtime.GOARCH == "arm64" && isRosettaInstalled() {
		// If on ARM64 and Rosetta 2 is installed, use the x86_64 version
		arch = "x86_64"
	}

	// Install Homebrew
	cmd = exec.Command("/bin/bash", "-c", "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)")
	cmd.Env = append(os.Environ(), "HOMEBREW_ARCH="+arch)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		close(done) // Stop the loading animation
		fmt.Println("\rError installing Homebrew:", err)
		return
	}

	// Stop the loading animation
	close(done)
	fmt.Printf("\r✅ Homebrew installed successfully.\n")
}

func (h *HomebrewInstaller) Title() string {
	return "Homebrew"
}

func (h *HomebrewInstaller) IsAlreadyInstalled() bool {
	return homebrewInstalled()
}

func (h *HomebrewInstaller) Description() string {
	alreadyInstalledMsg := ""
	if homebrewInstalled() {
		alreadyInstalledMsg = " (Already installed)"
	}
	return fmt.Sprintf("This option will install latest version of Homebrew.%s", alreadyInstalledMsg)
}

func homebrewInstalled() bool {
	_, err := exec.LookPath("brew")
	return err == nil
}

func isRosettaInstalled() bool {
	cmd := exec.Command("arch")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	return err != nil && strings.Contains(err.Error(), "exec format error")
}
