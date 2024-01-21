package install

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

type DotnetInstaller struct{}

func (d *DotnetInstaller) Install() {
	fmt.Println("Checking .NET SDK installation...")

	if d.IsAlreadyInstalled() {
		return
	}

	fmt.Println(".NET SDK is not installed. Installing...")

	// Display an animated hourglass loading indicator
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				fmt.Printf("\r⌛ Installing .NET SDK %s", loadingAnimation())
				time.Sleep(300 * time.Millisecond)
			}
		}
	}()

	// Install .NET SDK using Homebrew
	cmd := exec.Command("brew", "install", "dotnet-sdk")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	// Stop the loading animation
	close(done)

	if err != nil {
		fmt.Printf("\r❌ Error installing .NET SDK: %v\n", err)
		return
	}

	fmt.Printf("\r✅ .NET SDK installed successfully.\n")
}

func (d *DotnetInstaller) Title() string {
	return ".NET SDK"
}

func (d *DotnetInstaller) IsAlreadyInstalled() bool {
	_, err := exec.LookPath("dotnet")
	return err == nil
}

func (d *DotnetInstaller) Description() string {
	alreadyInstalled := ""
	if d.IsAlreadyInstalled() {
		alreadyInstalled = " (Already installed)"
	}
	return fmt.Sprintf("This option will install latest version of .NET SDK.%s", alreadyInstalled)
}
