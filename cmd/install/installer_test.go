package install

import (
"os"
"path/filepath"
"reflect"
"testing"
)

// --- Tests from main (parseShellList, shellConfigPath) ---

func TestParseShellList(t *testing.T) {
content := `
# List of valid login shells
/bin/bash
/bin/zsh

  /opt/homebrew/bin/zsh
`

got := parseShellList(content)
want := []string{"/bin/bash", "/bin/zsh", "/opt/homebrew/bin/zsh"}
if !reflect.DeepEqual(got, want) {
t.Fatalf("parseShellList() = %v, want %v", got, want)
}
}

func TestShellConfigPath(t *testing.T) {
home := "/Users/test"

tests := []struct {
name    string
shell   string
want    string
wantErr bool
}{
{name: "zsh", shell: "/bin/zsh", want: "/Users/test/.zshrc"},
{name: "bash", shell: "/bin/bash", want: "/Users/test/.bashrc"},
{name: "unsupported", shell: "/bin/fish", wantErr: true},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
got, err := shellConfigPath(home, tt.shell)
if tt.wantErr {
if err == nil {
t.Fatalf("expected error for shell %q", tt.shell)
}
return
}
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if got != tt.want {
t.Fatalf("shellConfigPath() = %q, want %q", got, tt.want)
}
})
}
}

// --- Registry unit tests (ticket-028) ---

func TestGetTitles_NonEmpty(t *testing.T) {
titles := GetTitles()
if len(titles) == 0 {
t.Fatal("GetTitles: expected at least one title")
}
}

func TestGetTitles_NoBlanks(t *testing.T) {
for _, title := range GetTitles() {
if title == "" {
t.Error("GetTitles: found blank title in result")
}
}
}

func TestGetInstallerByTitle_KnownInstaller(t *testing.T) {
known := []string{"Homebrew", "Git", "Go", "NVM (NodeJS LTS)", "Docker"}
for _, name := range known {
t.Run(name, func(t *testing.T) {
inst := GetInstallerByTitle(name)
if inst == nil {
t.Fatalf("GetInstallerByTitle(%q): expected non-nil", name)
}
if inst.Title() != name {
t.Errorf("Title() = %q, want %q", inst.Title(), name)
}
})
}
}

func TestGetInstallerByTitle_Unknown(t *testing.T) {
if inst := GetInstallerByTitle("__unknown__"); inst != nil {
t.Errorf("expected nil for unknown title, got %T", inst)
}
}

func TestGetDescriptionByTitle_KnownInstaller(t *testing.T) {
for _, title := range GetTitles() {
t.Run(title, func(t *testing.T) {
desc := GetDescriptionByTitle(title)
if desc == "" {
t.Errorf("GetDescriptionByTitle(%q): expected non-empty description", title)
}
})
}
}

func TestGetDescriptionByTitle_Unknown(t *testing.T) {
if desc := GetDescriptionByTitle("__unknown__"); desc != "" {
t.Errorf("expected empty string for unknown title, got %q", desc)
}
}

func TestGetAlreadyInstalledTools_ValidTitles(t *testing.T) {
for _, tool := range GetAlreadyInstalledTools() {
if GetInstallerByTitle(tool) == nil {
t.Errorf("GetAlreadyInstalledTools returned unregistered tool %q", tool)
}
}
}

func TestInstallerRegistry_AllTitlesAndDescriptions(t *testing.T) {
for _, inst := range installers {
t.Run(inst.Title(), func(t *testing.T) {
if inst.Title() == "" {
t.Error("installer has empty title")
}
if inst.Description() == "" {
t.Errorf("installer %q has empty description", inst.Title())
}
})
}
}

func TestGetTitles_MatchesRegistry(t *testing.T) {
titles := GetTitles()
if len(titles) != len(installers) {
t.Errorf("GetTitles() returned %d titles, registry has %d installers", len(titles), len(installers))
}
}

// --- Integration tests with mock PATH (ticket-029) ---

// lookPathInstallers lists installers whose detection is purely PATH-based
// (exec.LookPath). Installers using filesystem checks (iTerm2, Docker) or
// env-var checks (NVM) are excluded because they cannot be fully mocked via
// PATH manipulation alone.
var lookPathInstallers = []string{"Homebrew", "Git", "Go", "Visual Studio Code", ".NET SDK", "Python", "Zsh"}

// TestDetection_WithEmptyPATH verifies PATH-based IsAlreadyInstalled returns
// false when no binaries are on PATH.
func TestDetection_WithEmptyPATH(t *testing.T) {
origPATH := os.Getenv("PATH")
origNVM := os.Getenv("NVM_DIR")
t.Cleanup(func() {
os.Setenv("PATH", origPATH)
if origNVM == "" {
os.Unsetenv("NVM_DIR")
} else {
os.Setenv("NVM_DIR", origNVM)
}
})

empty := t.TempDir()
os.Setenv("PATH", empty)
os.Unsetenv("NVM_DIR")

for _, title := range lookPathInstallers {
t.Run(title, func(t *testing.T) {
inst := GetInstallerByTitle(title)
if inst == nil {
t.Skipf("installer %q not registered", title)
}
if inst.IsAlreadyInstalled() {
t.Errorf("IsAlreadyInstalled() = true with empty PATH for %q", title)
}
})
}
}

// TestDetection_WithFakeBrew verifies HomebrewInstaller detects a fake brew binary.
func TestDetection_WithFakeBrew(t *testing.T) {
orig := os.Getenv("PATH")
t.Cleanup(func() { os.Setenv("PATH", orig) })

dir := t.TempDir()
brewPath := filepath.Join(dir, "brew")
if err := os.WriteFile(brewPath, []byte("#!/bin/sh\necho fake\n"), 0755); err != nil {
t.Fatalf("WriteFile: %v", err)
}
os.Setenv("PATH", dir)

inst := GetInstallerByTitle("Homebrew")
if inst == nil {
t.Fatal("Homebrew installer not found")
}
if !inst.IsAlreadyInstalled() {
t.Error("HomebrewInstaller.IsAlreadyInstalled() = false with fake brew on PATH")
}
}

// TestDetection_WithFakeGit verifies GitInstaller detects a fake git binary.
func TestDetection_WithFakeGit(t *testing.T) {
orig := os.Getenv("PATH")
t.Cleanup(func() { os.Setenv("PATH", orig) })

dir := t.TempDir()
gitPath := filepath.Join(dir, "git")
if err := os.WriteFile(gitPath, []byte("#!/bin/sh\necho fake\n"), 0755); err != nil {
t.Fatalf("WriteFile: %v", err)
}
os.Setenv("PATH", dir)

inst := GetInstallerByTitle("Git")
if inst == nil {
t.Fatal("Git installer not found")
}
if !inst.IsAlreadyInstalled() {
t.Error("GitInstaller.IsAlreadyInstalled() = false with fake git on PATH")
}
}

// TestDetection_WithFakeGo verifies GoInstaller detects a fake go binary.
func TestDetection_WithFakeGo(t *testing.T) {
orig := os.Getenv("PATH")
t.Cleanup(func() { os.Setenv("PATH", orig) })

dir := t.TempDir()
goPath := filepath.Join(dir, "go")
if err := os.WriteFile(goPath, []byte("#!/bin/sh\necho fake\n"), 0755); err != nil {
t.Fatalf("WriteFile: %v", err)
}
os.Setenv("PATH", dir)

inst := GetInstallerByTitle("Go")
if inst == nil {
t.Fatal("Go installer not found")
}
if !inst.IsAlreadyInstalled() {
t.Error("GoInstaller.IsAlreadyInstalled() = false with fake go on PATH")
}
}
