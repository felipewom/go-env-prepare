package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildPlanAutoDiscoverySuccess(t *testing.T) {
	tmp := t.TempDir()
	manifest := "apiVersion: v1\nprofile: backend\n"
	if err := os.WriteFile(filepath.Join(tmp, "prepare.yaml"), []byte(manifest), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	prevWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	defer func() { _ = os.Chdir(prevWD) }()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	plan, resolvedManifest, err := buildPlan(&dynamicFlags{})
	if err != nil {
		t.Fatalf("buildPlan error: %v", err)
	}
	if resolvedManifest.Profile != "backend" {
		t.Fatalf("expected discovered manifest profile backend, got %q", resolvedManifest.Profile)
	}
	if len(plan.Steps) == 0 {
		t.Fatal("expected non-empty execution plan")
	}
}

func TestBuildPlanNoManifestFallback(t *testing.T) {
	tmp := t.TempDir()
	prevWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	defer func() { _ = os.Chdir(prevWD) }()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	plan, resolvedManifest, err := buildPlan(&dynamicFlags{})
	if err != nil {
		t.Fatalf("buildPlan error: %v", err)
	}
	if resolvedManifest.Profile != "" {
		t.Fatalf("expected empty manifest profile fallback, got %q", resolvedManifest.Profile)
	}
	if len(plan.Steps) == 0 {
		t.Fatal("expected fallback builtin plan")
	}
}

func TestBuildPlanInvalidDiscoveredManifestReturnsError(t *testing.T) {
	tmp := t.TempDir()
	invalid := "apiVersion: v1\nprofiles:\n  bad: [\n"
	if err := os.WriteFile(filepath.Join(tmp, "prepare.yaml"), []byte(invalid), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	prevWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	defer func() { _ = os.Chdir(prevWD) }()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	_, _, err = buildPlan(&dynamicFlags{})
	if err == nil {
		t.Fatal("expected error for invalid discovered manifest")
	}
	if !strings.Contains(err.Error(), "unsupported manifest syntax") {
		t.Fatalf("unexpected error: %v", err)
	}
}
