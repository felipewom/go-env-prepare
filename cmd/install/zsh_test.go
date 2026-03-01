package install

import (
	"os/user"
	"path/filepath"
	"testing"
)

func TestResolveZshCustomPath_Default(t *testing.T) {
	t.Setenv("ZSH_CUSTOM", "")

	usr, err := user.Current()
	if err != nil {
		t.Fatalf("user.Current() error: %v", err)
	}

	got, err := resolveZshCustomPath()
	if err != nil {
		t.Fatalf("resolveZshCustomPath() error: %v", err)
	}

	want := filepath.Join(usr.HomeDir, ".oh-my-zsh", "custom")
	if got != want {
		t.Fatalf("resolveZshCustomPath() = %q, want %q", got, want)
	}
}

func TestResolveZshCustomPath_ExpandEnv(t *testing.T) {
	t.Setenv("CUSTOM_SUFFIX", "custom-dir")
	t.Setenv("ZSH_CUSTOM", "/tmp/$CUSTOM_SUFFIX")

	got, err := resolveZshCustomPath()
	if err != nil {
		t.Fatalf("resolveZshCustomPath() error: %v", err)
	}

	want := "/tmp/custom-dir"
	if got != want {
		t.Fatalf("resolveZshCustomPath() = %q, want %q", got, want)
	}
}

func TestResolveZshCustomPath_ExpandTilde(t *testing.T) {
	t.Setenv("ZSH_CUSTOM", "~/my-custom")

	usr, err := user.Current()
	if err != nil {
		t.Fatalf("user.Current() error: %v", err)
	}

	got, err := resolveZshCustomPath()
	if err != nil {
		t.Fatalf("resolveZshCustomPath() error: %v", err)
	}

	want := filepath.Join(usr.HomeDir, "my-custom")
	if got != want {
		t.Fatalf("resolveZshCustomPath() = %q, want %q", got, want)
	}
}
