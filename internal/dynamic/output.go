package dynamic

import (
	"encoding/json"
	"fmt"
	"os"
)

func PrintJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func PrintPlanHuman(plan Plan) {
	fmt.Printf("Execution plan (%d steps):\n", len(plan.Steps))
	for _, step := range plan.Steps {
		fmt.Printf("%d. %s (%s)\n", step.Order, step.Tool.Title, step.Tool.ID)
	}
}

func PrintExecutionHuman(result ExecutionResult) {
	fmt.Printf("Run summary: %d steps, dry-run=%v\n", len(result.Steps), result.DryRun)
	for _, step := range result.Steps {
		status := "ok"
		if !step.Success {
			status = "failed"
		}
		if step.Action == "skip" {
			status = "skipped"
		}
		fmt.Printf("- %s: %s", step.ToolID, status)
		if step.Reason != "" {
			fmt.Printf(" (%s)", step.Reason)
		}
		if step.Error != "" {
			fmt.Printf(" error=%s", step.Error)
		}
		fmt.Printf("\n")
	}
}
