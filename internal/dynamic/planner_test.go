package dynamic

import "testing"

func TestBuildPlanResolvesDependenciesFirst(t *testing.T) {
	plan, err := BuildPlan([]string{"docker"}, BuiltinCatalog())
	if err != nil {
		t.Fatalf("BuildPlan error: %v", err)
	}
	if len(plan.Steps) < 2 {
		t.Fatalf("expected at least 2 steps, got %d", len(plan.Steps))
	}
	if plan.Steps[0].Tool.ID != "homebrew" {
		t.Fatalf("expected homebrew first, got %s", plan.Steps[0].Tool.ID)
	}
	if plan.Steps[len(plan.Steps)-1].Tool.ID != "docker" {
		t.Fatalf("expected docker last, got %s", plan.Steps[len(plan.Steps)-1].Tool.ID)
	}
}

func TestBuildPlanCycleDetection(t *testing.T) {
	catalog := map[string]ToolSpec{
		"a": {ID: "a", Dependencies: []string{"b"}},
		"b": {ID: "b", Dependencies: []string{"a"}},
	}
	_, err := BuildPlan([]string{"a"}, catalog)
	if err == nil {
		t.Fatal("expected cycle detection error")
	}
}
