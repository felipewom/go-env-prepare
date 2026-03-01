package dynamic

import "time"

type Manifest struct {
	APIVersion string             `json:"apiVersion"`
	Profile    string             `json:"profile,omitempty"`
	Tools      []string           `json:"tools,omitempty"`
	Profiles   map[string]Profile `json:"profiles,omitempty"`
}

type Profile struct {
	Extends []string `json:"extends,omitempty"`
	Tools   []string `json:"tools,omitempty"`
}

type ToolSpec struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Dependencies []string `json:"dependencies,omitempty"`
	Install      Command  `json:"install"`
	Check        Check    `json:"check"`
	Version      string   `json:"version,omitempty"`
	Source       string   `json:"source,omitempty"`
}

type Command struct {
	Name  string   `json:"name"`
	Args  []string `json:"args,omitempty"`
	Shell string   `json:"shell,omitempty"`
}

type Check struct {
	Binary     string `json:"binary,omitempty"`
	PathExists string `json:"pathExists,omitempty"`
}

type Plan struct {
	CreatedAt time.Time  `json:"createdAt"`
	Steps     []PlanStep `json:"steps"`
}

type PlanStep struct {
	Order int      `json:"order"`
	Tool  ToolSpec `json:"tool"`
}

type ExecutionResult struct {
	StartedAt time.Time       `json:"startedAt"`
	EndedAt   time.Time       `json:"endedAt"`
	DryRun    bool            `json:"dryRun"`
	Steps     []ExecutionStep `json:"steps"`
}

type ExecutionStep struct {
	ToolID     string `json:"toolId"`
	Action     string `json:"action"`
	Reason     string `json:"reason,omitempty"`
	Success    bool   `json:"success"`
	Error      string `json:"error,omitempty"`
	DurationMs int64  `json:"durationMs"`
}

type State struct {
	Completed map[string]bool `json:"completed"`
}

type Lockfile struct {
	Version     int          `json:"version"`
	GeneratedAt time.Time    `json:"generatedAt"`
	Tools       []LockedTool `json:"tools"`
}

type LockedTool struct {
	ID      string `json:"id"`
	Version string `json:"version,omitempty"`
	Source  string `json:"source,omitempty"`
}
