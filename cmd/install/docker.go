package install

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

type DockerInstaller struct{}

func (d *DockerInstaller) Install() {
	fmt.Println("Checking Docker installation...")

	if d.IsAlreadyInstalled() {
		return
	}

	fmt.Println("Docker is not installed. Installing...")

	// Display an animated hourglass loading indicator
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				fmt.Printf("\r⌛ Installing Docker %s", loadingAnimation())
				time.Sleep(300 * time.Millisecond)
			}
		}
	}()

	// Install Docker Desktop using Homebrew cask
	cmd := exec.Command("brew", "install", "--cask", "docker")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	// Stop the loading animation
	close(done)

	if err != nil {
		fmt.Printf("\r❌ Error installing Docker: %v\n", err)
		return
	}

	fmt.Printf("\r✅ Docker installed successfully.\n")

	// Perform post-installation guidance
	d.configureDocker()
}

func (d *DockerInstaller) Title() string {
	return "Docker"
}

func (d *DockerInstaller) IsAlreadyInstalled() bool {
	_, appErr := os.Stat("/Applications/Docker.app")
	_, cliErr := exec.LookPath("docker")
	return appErr == nil || cliErr == nil
}

func (d *DockerInstaller) Description() string {
	alreadyInstalled := ""
	if d.IsAlreadyInstalled() {
		alreadyInstalled = " (Already installed)"
	}
	return fmt.Sprintf("This option will install latest version of Docker.%s", alreadyInstalled)
}

func (d *DockerInstaller) configureDocker() {
	fmt.Println("Docker Desktop installed. Open Docker.app once to finish initialization.")
}
