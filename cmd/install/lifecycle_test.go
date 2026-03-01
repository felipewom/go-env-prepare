package install

import "testing"

type stubInstaller struct {
	installed bool
	installN  int
	checkN    int
	applyN    int
}

func (s *stubInstaller) Install() {
	s.installN++
	s.installed = true
}

func (s *stubInstaller) IsAlreadyInstalled() bool {
	return s.installed
}

func (s *stubInstaller) Title() string {
	return "stub"
}

func (s *stubInstaller) Description() string {
	return "stub"
}

func (s *stubInstaller) Plan() string {
	return "plan"
}

func (s *stubInstaller) Check() bool {
	s.checkN++
	return s.installed
}

func (s *stubInstaller) Apply() error {
	s.applyN++
	s.installed = true
	return nil
}

func TestRunInstaller_UsesLifecycleWhenAvailable(t *testing.T) {
	s := &stubInstaller{}

	RunInstaller(s)

	if s.checkN != 1 {
		t.Fatalf("expected Check to be called once, got %d", s.checkN)
	}
	if s.applyN != 1 {
		t.Fatalf("expected Apply to be called once, got %d", s.applyN)
	}
	if s.installN != 0 {
		t.Fatalf("expected Install not to be called, got %d", s.installN)
	}
}

func TestIsInstalled_UsesLifecycleCheck(t *testing.T) {
	s := &stubInstaller{installed: true}

	if !IsInstalled(s) {
		t.Fatal("expected installer to be detected as installed")
	}
	if s.checkN != 1 {
		t.Fatalf("expected Check to be called once, got %d", s.checkN)
	}
}
