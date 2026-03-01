package install

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

type ZshInstaller struct{}

const zshrcManagedBlock = `
# >>> go-env-prepare zsh >>>
# Plugins installed by go-env-prepare
if typeset -p plugins >/dev/null 2>&1; then
  plugins+=(docker nvm go zsh-autosuggestions zsh-syntax-highlighting)
  typeset -U plugins
fi

alias zc='code ~/.zshrc'
alias rr='source ~/.zshrc && clear && echo "Zsh configuration reloaded."'
# <<< go-env-prepare zsh <<<
`

func (z *ZshInstaller) Install() {
	if z.Check() {
		return
	}

	if err := z.Apply(); err != nil {
		fmt.Printf("\r❌ Error installing Zsh: %v\n", err)
	}
}

func (z *ZshInstaller) Plan() string {
	return "Install Zsh, set it as default shell, install Oh My Zsh/plugins, and update .zshrc."
}

func (z *ZshInstaller) Check() bool {
	return z.IsAlreadyInstalled() && z.isConfigured()
}

func (z *ZshInstaller) Apply() error {
	fmt.Println("Checking Zsh installation...")

	if !z.IsAlreadyInstalled() {
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

		cmd := exec.Command("brew", "install", "zsh")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		close(done)
		if err != nil {
			return err
		}

		fmt.Printf("\r✅ Zsh installed successfully.\n")
	} else {
		fmt.Printf("\r✅ Zsh is already installed.\n")
	}

	z.setDefaultShell()
	z.installOhMyZsh()
	z.postInstall()
	return nil
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

// IsZshInstalled checks if Zsh is already installed.
func IsZshInstalled() bool {
	_, err := exec.LookPath("zsh")
	return err == nil
}

func (z *ZshInstaller) isConfigured() bool {
	usr, err := user.Current()
	if err != nil {
		return false
	}

	zshrcPath := filepath.Join(usr.HomeDir, ".zshrc")
	content, err := os.ReadFile(zshrcPath)
	if err != nil {
		return false
	}
	if !strings.Contains(string(content), "# >>> go-env-prepare zsh >>>") {
		return false
	}

	customPath, err := resolveZshCustomPath()
	if err != nil {
		return false
	}

	if _, err := os.Stat(filepath.Join(customPath, "plugins", "zsh-autosuggestions")); err != nil {
		return false
	}
	if _, err := os.Stat(filepath.Join(customPath, "plugins", "zsh-syntax-highlighting")); err != nil {
		return false
	}

	return true
}

func (z *ZshInstaller) setDefaultShell() {
	fmt.Println("Setting Zsh as the default shell...")

	path, err := detectValidZshShellPath()
	if err != nil {
		fmt.Printf("\r⚠️ Skipping default shell change: %v\n", err)
		return
	}

	if os.Getenv("SHELL") == path {
		fmt.Printf("\r✅ Zsh is already the active shell.\n")
		return
	}

	cmd := exec.Command("chsh", "-s", path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Printf("\r❌ Error setting Zsh as the default shell: %v\n", err)
		return
	}

	fmt.Printf("\r✅ Zsh set as the default shell.\n")
}

func detectValidZshShellPath() (string, error) {
	zshPath, err := exec.LookPath("zsh")
	if err != nil {
		return "", errors.New("zsh binary not found")
	}

	if !shellListedInEtcShells(zshPath) {
		return "", fmt.Errorf("%s is not listed in /etc/shells", zshPath)
	}

	return zshPath, nil
}

func (z *ZshInstaller) installOhMyZsh() {
	fmt.Println("Installing Oh My Zsh...")

	usr, err := user.Current()
	if err != nil {
		fmt.Printf("\r❌ Error getting user's home directory: %v\n", err)
		return
	}

	ohMyZshDir := filepath.Join(usr.HomeDir, ".oh-my-zsh")
	if _, err := os.Stat(ohMyZshDir); err == nil {
		fmt.Printf("\r✅ Oh My Zsh is already installed.\n")
		return
	}

	cmd := exec.Command("sh", "-c", "$(curl -fsSL https://raw.github.com/ohmyzsh/ohmyzsh/master/tools/install.sh)")
	cmd.Env = append(os.Environ(), "RUNZSH=no", "CHSH=no", "KEEP_ZSHRC=yes")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Printf("\r❌ Error installing Oh My Zsh: %v\n", err)
		return
	}

	fmt.Printf("\r✅ Oh My Zsh installed successfully.\n")
}

func (z *ZshInstaller) postInstall() {
	fmt.Println("Performing post-installation tasks for Zsh...")
	z.installZshPlugins()
	z.updateZshrc()
}

func (z *ZshInstaller) installZshPlugins() {
	fmt.Println("Installing Zsh plugins...")

	customPath, err := resolveZshCustomPath()
	if err != nil {
		fmt.Printf("\r❌ Error resolving ZSH_CUSTOM path: %v\n", err)
		return
	}

	pluginsDir := filepath.Join(customPath, "plugins")
	if err := os.MkdirAll(pluginsDir, 0o755); err != nil {
		fmt.Printf("\r❌ Error creating plugin directory: %v\n", err)
		return
	}

	z.clonePluginIfMissing("https://github.com/zsh-users/zsh-autosuggestions.git", filepath.Join(pluginsDir, "zsh-autosuggestions"))
	z.clonePluginIfMissing("https://github.com/zsh-users/zsh-syntax-highlighting.git", filepath.Join(pluginsDir, "zsh-syntax-highlighting"))

	fmt.Printf("\r✅ Zsh plugins installed successfully.\n")
}

func (z *ZshInstaller) clonePluginIfMissing(repo, destination string) {
	if _, err := os.Stat(destination); err == nil {
		fmt.Printf("\r✅ Plugin already installed: %s\n", filepath.Base(destination))
		return
	}

	cmd := exec.Command("git", "clone", repo, destination)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("\r❌ Error installing %s: %v\n", filepath.Base(destination), err)
	}
}

func resolveZshCustomPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	customPath := os.Getenv("ZSH_CUSTOM")
	if customPath == "" {
		customPath = filepath.Join(usr.HomeDir, ".oh-my-zsh", "custom")
	}

	customPath = os.ExpandEnv(customPath)
	if strings.HasPrefix(customPath, "~/") {
		customPath = filepath.Join(usr.HomeDir, customPath[2:])
	}
	if !filepath.IsAbs(customPath) {
		customPath = filepath.Join(usr.HomeDir, customPath)
	}

	return filepath.Clean(customPath), nil
}

func (z *ZshInstaller) updateZshrc() {
	fmt.Println("Updating .zshrc file...")

	usr, err := user.Current()
	if err != nil {
		fmt.Printf("\r❌ Error getting user's home directory: %v\n", err)
		return
	}

	zshrcPath := filepath.Join(usr.HomeDir, ".zshrc")
	content, err := os.ReadFile(zshrcPath)
	if err != nil && !os.IsNotExist(err) {
		fmt.Printf("\r❌ Error reading .zshrc file: %v\n", err)
		return
	}

	if strings.Contains(string(content), "# >>> go-env-prepare zsh >>>") {
		fmt.Printf("\r✅ .zshrc already contains go-env-prepare snippet.\n")
		return
	}

	f, err := os.OpenFile(zshrcPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		fmt.Printf("\r❌ Error opening .zshrc file: %v\n", err)
		return
	}
	defer f.Close()

	_, err = f.WriteString(zshrcManagedBlock)
	if err != nil {
		fmt.Printf("\r❌ Error writing to .zshrc file: %v\n", err)
		return
	}

	fmt.Printf("\r✅ .zshrc file updated successfully.\n")
}
