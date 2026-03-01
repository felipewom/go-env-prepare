package dynamic

import (
	"sort"
	"time"
)

func BuildLockfile(plan Plan) Lockfile {
	tools := make([]LockedTool, 0, len(plan.Steps))
	seen := map[string]bool{}
	for _, step := range plan.Steps {
		if seen[step.Tool.ID] {
			continue
		}
		seen[step.Tool.ID] = true
		tools = append(tools, LockedTool{
			ID:      step.Tool.ID,
			Version: step.Tool.Version,
			Source:  step.Tool.Source,
		})
	}
	sort.Slice(tools, func(i, j int) bool { return tools[i].ID < tools[j].ID })
	return Lockfile{
		Version:     1,
		GeneratedAt: time.Now(),
		Tools:       tools,
	}
}
