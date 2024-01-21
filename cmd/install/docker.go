package install

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
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
				fmt.Printf("\r⌛ Installing Docker %s \n", loadingAnimation())
				time.Sleep(300 * time.Millisecond)
			}
		}
	}()

	// Install Docker using Homebrew
	cmd := exec.Command("brew", "install", "docker")
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

	// Perform post-installation configurations
	d.configureDocker()
}

func (d *DockerInstaller) Title() string {
	return "Docker"
}

func (d *DockerInstaller) IsAlreadyInstalled() bool {
	_, err := exec.LookPath("docker")
	return err == nil
}

func (d *DockerInstaller) Description() string {
	alreadyInstalled := ""
	if d.IsAlreadyInstalled() {
		alreadyInstalled = " (Already installed)"
	}
	return fmt.Sprintf("This option will install latest version of Docker.%s", alreadyInstalled)
}

func (d *DockerInstaller) configureDocker() {
	fmt.Println("Performing post-installation configurations for Docker...")

	// Get the user's username
	currentUser, err := user.Current()
	if err != nil {
		fmt.Printf("\r❌ Error getting user information: %v\n", err)
		return
	}

	// Add the user to the 'docker' group using dscl
	cmd := exec.Command("dscl", ".", "-append", "/Groups/docker", "GroupMembership", currentUser.Username)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()

	if err != nil {
		fmt.Printf("\r❌ Error configuring Docker: %v\n", err)
		return
	}

	fmt.Printf("\r✅ Docker configured successfully. Please restart your shell.\n")
}
