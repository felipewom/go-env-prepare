package dynamic

import (
	"runtime"
	"testing"
)

func TestExecutorDryRunDoesNotExecuteCommand(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("preflight enforces darwin")
	}
	executor := NewExecutor()
	called := false
	executor.checkTool = func(c Check) bool { return false }
	executor.runCommand = func(cmd Command) error {
		called = true
		return nil
	}

	plan := Plan{Steps: []PlanStep{{Order: 1, Tool: ToolSpec{ID: "git", Install: Command{Name: "echo"}}}}}
	result, err := executor.Run(plan, ExecOptions{DryRun: true})
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	if called {
		t.Fatal("expected command runner not to be called in dry-run")
	}
	if len(result.Steps) != 1 || result.Steps[0].Reason != "dry_run" {
		t.Fatalf("unexpected dry-run result: %#v", result.Steps)
	}
}
