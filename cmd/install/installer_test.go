package install

import (
	"reflect"
	"testing"
)

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
