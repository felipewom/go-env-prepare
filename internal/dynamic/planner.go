package dynamic

import (
	"fmt"
	"sort"
	"time"
)

func BuildPlan(toolIDs []string, catalog map[string]ToolSpec) (Plan, error) {
	expanded := map[string]bool{}
	visiting := map[string]bool{}
	for _, id := range toolIDs {
		if err := includeWithDependencies(id, catalog, expanded, visiting); err != nil {
			return Plan{}, err
		}
	}

	inDegree := map[string]int{}
	graph := map[string][]string{}
	expandedIDs := make([]string, 0, len(expanded))
	for id := range expanded {
		inDegree[id] = 0
		expandedIDs = append(expandedIDs, id)
	}
	sort.Strings(expandedIDs)
	for _, id := range expandedIDs {
		for _, dep := range catalog[id].Dependencies {
			if !expanded[dep] {
				continue
			}
			graph[dep] = append(graph[dep], id)
			inDegree[id]++
		}
	}

	queue := []string{}
	for _, requested := range toolIDs {
		if inDegree[requested] == 0 && expanded[requested] {
			queue = appendUnique(queue, requested)
		}
	}
	for _, id := range expandedIDs {
		if inDegree[id] == 0 {
			queue = appendUnique(queue, id)
		}
	}

	ordered := []string{}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		ordered = append(ordered, node)
		sort.Strings(graph[node])
		for _, next := range graph[node] {
			inDegree[next]--
			if inDegree[next] == 0 {
				queue = append(queue, next)
			}
		}
		sort.Strings(queue)
	}

	if len(ordered) != len(expanded) {
		return Plan{}, fmt.Errorf("dependency cycle detected while building execution plan")
	}

	steps := make([]PlanStep, 0, len(ordered))
	for idx, id := range ordered {
		steps = append(steps, PlanStep{Order: idx + 1, Tool: catalog[id]})
	}

	return Plan{CreatedAt: time.Now(), Steps: steps}, nil
}

func includeWithDependencies(id string, catalog map[string]ToolSpec, out map[string]bool, visiting map[string]bool) error {
	if out[id] {
		return nil
	}
	if visiting[id] {
		return fmt.Errorf("dependency cycle detected at tool %q", id)
	}
	spec, ok := catalog[id]
	if !ok {
		return fmt.Errorf("unknown tool %q", id)
	}
	visiting[id] = true
	for _, dep := range spec.Dependencies {
		if err := includeWithDependencies(dep, catalog, out, visiting); err != nil {
			return err
		}
	}
	delete(visiting, id)
	out[id] = true
	return nil
}

func appendUnique(values []string, value string) []string {
	for _, v := range values {
		if v == value {
			return values
		}
	}
	return append(values, value)
}
