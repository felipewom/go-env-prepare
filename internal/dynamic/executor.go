package dynamic

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"
)

type ExecOptions struct {
	DryRun    bool
	Resume    bool
	StatePath string
}

type commandRunner func(cmd Command) error

type checker func(c Check) bool

type Executor struct {
	runCommand commandRunner
	checkTool  checker
}

func NewExecutor() *Executor {
	return &Executor{
		runCommand: defaultCommandRunner,
		checkTool:  defaultChecker,
	}
}

func (e *Executor) Run(plan Plan, opts ExecOptions) (ExecutionResult, error) {
	if err := preflight(); err != nil {
		return ExecutionResult{}, err
	}
	if opts.StatePath == "" {
		opts.StatePath = ".prepare.state.json"
	}

	state := State{Completed: map[string]bool{}}
	var err error
	if opts.Resume {
		state, err = LoadState(opts.StatePath)
		if err != nil {
			return ExecutionResult{}, fmt.Errorf("load state: %w", err)
		}
	}

	result := ExecutionResult{StartedAt: time.Now(), DryRun: opts.DryRun, Steps: []ExecutionStep{}}
	for _, step := range plan.Steps {
		stepStart := time.Now()
		execStep := ExecutionStep{ToolID: step.Tool.ID, Success: true}

		if state.Completed[step.Tool.ID] {
			execStep.Action = "skip"
			execStep.Reason = "already_completed"
			execStep.DurationMs = time.Since(stepStart).Milliseconds()
			result.Steps = append(result.Steps, execStep)
			continue
		}
		if e.checkTool(step.Tool.Check) {
			execStep.Action = "skip"
			execStep.Reason = "already_installed"
			execStep.DurationMs = time.Since(stepStart).Milliseconds()
			result.Steps = append(result.Steps, execStep)
			state.Completed[step.Tool.ID] = true
			if opts.Resume {
				if err := SaveState(opts.StatePath, state); err != nil {
					return result, fmt.Errorf("save state: %w", err)
				}
			}
			continue
		}
		if opts.DryRun {
			execStep.Action = "install"
			execStep.Reason = "dry_run"
			execStep.DurationMs = time.Since(stepStart).Milliseconds()
			result.Steps = append(result.Steps, execStep)
			continue
		}

		execStep.Action = "install"
		if err := e.runCommand(step.Tool.Install); err != nil {
			execStep.Success = false
			execStep.Error = err.Error()
			execStep.DurationMs = time.Since(stepStart).Milliseconds()
			result.Steps = append(result.Steps, execStep)
			result.EndedAt = time.Now()
			return result, fmt.Errorf("install %s: %w", step.Tool.ID, err)
		}

		state.Completed[step.Tool.ID] = true
		if opts.Resume {
			if err := SaveState(opts.StatePath, state); err != nil {
				return result, fmt.Errorf("save state: %w", err)
			}
		}

		execStep.DurationMs = time.Since(stepStart).Milliseconds()
		result.Steps = append(result.Steps, execStep)
	}

	result.EndedAt = time.Now()
	return result, nil
}

func preflight() error {
	if runtime.GOOS != "darwin" {
		return fmt.Errorf("unsupported OS: %s (macOS required)", runtime.GOOS)
	}
	return nil
}

func defaultCommandRunner(c Command) error {
	var cmd *exec.Cmd
	if c.Shell != "" {
		cmd = exec.Command(c.Name, c.Args...)
	} else {
		cmd = exec.Command(c.Name, c.Args...)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func defaultChecker(c Check) bool {
	if c.Binary != "" {
		_, err := exec.LookPath(c.Binary)
		if err == nil {
			return true
		}
	}
	if c.PathExists != "" {
		if _, err := os.Stat(c.PathExists); err == nil {
			return true
		}
	}
	return false
}
