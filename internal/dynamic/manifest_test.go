package dynamic

import "testing"

func TestResolveToolsWithProfileInheritance(t *testing.T) {
	m := Manifest{
		APIVersion: "v1",
		Profile:    "custom",
		Profiles: map[string]Profile{
			"custom": {
				Extends: []string{"backend"},
				Tools:   []string{"nodejs"},
			},
		},
	}
	tools, err := ResolveTools(m, "", BuiltinCatalog(), BuiltinProfiles())
	if err != nil {
		t.Fatalf("ResolveTools error: %v", err)
	}

	expected := []string{"homebrew", "git", "zsh", "vscode", "go", "python", "docker", "nodejs"}
	for _, want := range expected {
		found := false
		for _, got := range tools {
			if got == want {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("expected tool %q in resolved list %#v", want, tools)
		}
	}
}

func TestValidateManifestUnknownTool(t *testing.T) {
	m := Manifest{APIVersion: "v1", Tools: []string{"not-real"}}
	err := ValidateManifest(m, BuiltinCatalog(), BuiltinProfiles())
	if err == nil {
		t.Fatal("expected validation error for unknown tool")
	}
}
