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

type PythonInstaller struct{}

func (p *PythonInstaller) Install() {
	fmt.Println("Checking Python installation...")

	if p.IsAlreadyInstalled() {
		return
	}

	fmt.Println("Python is not installed. Installing...")

	// Display an animated hourglass loading indicator
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				fmt.Printf("\r⌛ Installing Python %s", loadingAnimation())
				time.Sleep(300 * time.Millisecond)
			}
		}
	}()

	// Install Python using Homebrew
	cmd := exec.Command("brew", "install", "python")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	// Stop the loading animation
	close(done)

	if err != nil {
		fmt.Printf("\r❌ Error installing Python: %v\n", err)
		return
	}

	fmt.Printf("\r✅ Python installed successfully.\n")

	// Install Pyenv
	p.installPyenv()

	// Set Pyenv environment variables in shell configuration
	p.setPyenvVariables()

	// Install the latest version of Python
	p.installLatestPython()
}

func (p *PythonInstaller) Title() string {
	return "Python"
}

func (p *PythonInstaller) Description() string {
	alreadyInstalled := ""
	if p.IsAlreadyInstalled() {
		alreadyInstalled = " (Already installed)"
	}
	return fmt.Sprintf("This option will install Python and Pyenv.%s", alreadyInstalled)
}

func (p *PythonInstaller) IsAlreadyInstalled() bool {
	_, err := exec.LookPath("python")
	return err == nil
}

func (p *PythonInstaller) installPyenv() {
	fmt.Println("Installing Pyenv...")

	// Install Pyenv using Homebrew
	cmd := exec.Command("brew", "install", "pyenv")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		fmt.Printf("\r❌ Error installing Pyenv: %v\n", err)
		return
	}

	fmt.Printf("\r✅ Pyenv installed successfully.\n")
}

func (p *PythonInstaller) setPyenvVariables() {
	fmt.Println("Setting Pyenv environment variables...")

	// Get the user's home directory
	usr, err := user.Current()
	if err != nil {
		fmt.Printf("\r❌ Error getting user's home directory: %v\n", err)
		return
	}

	// Update the shell configuration file to set Pyenv variables
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

	// check if the shell configuration contains 'pyenv' already
	cmd := exec.Command("bash", "-c", fmt.Sprintf(`grep -q -F 'pyenv' %s`, configFile))
	err = cmd.Run()
	if err == nil {
		fmt.Printf("\r✅ Pyenv environment variables set successfully. Please restart your shell.\n")
		return
	}

	pyenvConfig := `

# Pyenv
export PYENV_ROOT="$HOME/.pyenv"
[[ -d $PYENV_ROOT/bin ]] && export PATH="$PYENV_ROOT/bin:$PATH"
eval "$(pyenv init -)"`

	// Add Pyenv initialization to the shell configuration file
	cmd = exec.Command("bash", "-c", fmt.Sprintf(`echo '%s' >> %s`, pyenvConfig, configFile))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Printf("\r❌ Error setting Pyenv environment variables: %v\n", err)
		return
	}

	cmd = exec.Command("bash", "-c", fmt.Sprintf(`echo 'eval "$(pyenv init --path)"' >> %s`, configFile))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Printf("\r❌ Error setting Pyenv environment variables: %v\n", err)
		return
	}

	fmt.Printf("\r✅ Pyenv environment variables set successfully. Please restart your shell.\n")
}

func (p *PythonInstaller) installLatestPython() {
	fmt.Println("Installing the latest version of Python...")

	// reload shell configuration
	err := reloadShellConfiguration()
	if err != nil {
		fmt.Printf("\r❌ Error reloading shell configuration: %v\n", err)
		return
	}

	// Source Pyenv in the current shell
	cmd := exec.Command("bash", "-c", "pyenv install $(pyenv install --list | grep -v - | tail -1) && pyenv global $(pyenv install --list | grep -v - | tail -1)")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()

	if err != nil {
		fmt.Printf("\r❌ Error installing Python: %v\n", err)
		return
	}

	fmt.Printf("\r✅ Python installed successfully.\n")
}
