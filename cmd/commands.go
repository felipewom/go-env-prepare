package cmd

import (
	"encoding/json"
	"errors"
	"felipewom/go-env-prepare/cmd/install"
	"felipewom/go-env-prepare/internal/dynamic"
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

type Commands struct {
	RootCmd *cobra.Command
}

type dynamicFlags struct {
	ManifestPath string
	Profile      string
	OutputJSON   bool
	DryRun       bool
	Resume       bool
	StatePath    string
	LockfilePath string
}

func NewCommands() *Commands {
	rootCmd := &Commands{RootCmd: newInstallCmd()}
	rootCmd.RootCmd.AddCommand(newVersionCmd())
	rootCmd.RootCmd.AddCommand(newPlanCmd())
	rootCmd.RootCmd.AddCommand(newRunCmd())
	rootCmd.RootCmd.AddCommand(newLintCmd())
	rootCmd.RootCmd.AddCommand(newLockCmd())
	return rootCmd
}

func newInstallCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "prepare",
		Short: "Go Env Prepare is a CLI tool for preparing your development environment",
		Long:  `With Go Env Prepare you can easily prepare your development environment for Go, NodeJS, Docker, etc.`,
		Run: func(cmd *cobra.Command, args []string) {
			startPrompt()
		},
	}
}

func newPlanCmd() *cobra.Command {
	flags := &dynamicFlags{}
	cmd := &cobra.Command{
		Use:   "plan",
		Short: "Build an execution plan from builtin or manifest profiles",
		RunE: func(cmd *cobra.Command, args []string) error {
			plan, _, err := buildPlan(flags)
			if err != nil {
				return err
			}
			if flags.OutputJSON {
				return dynamic.PrintJSON(plan)
			}
			dynamic.PrintPlanHuman(plan)
			return nil
		},
	}
	bindDynamicFlags(cmd, flags)
	return cmd
}

func newRunCmd() *cobra.Command {
	flags := &dynamicFlags{}
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Execute profile plan with idempotency and optional dry-run",
		RunE: func(cmd *cobra.Command, args []string) error {
			plan, _, err := buildPlan(flags)
			if err != nil {
				return err
			}
			executor := dynamic.NewExecutor()
			result, runErr := executor.Run(plan, dynamic.ExecOptions{
				DryRun:    flags.DryRun,
				Resume:    flags.Resume,
				StatePath: flags.StatePath,
			})
			if flags.OutputJSON {
				if err := dynamic.PrintJSON(result); err != nil {
					return err
				}
			} else {
				dynamic.PrintExecutionHuman(result)
			}
			if runErr != nil {
				return runErr
			}
			if flags.LockfilePath != "" {
				lock := dynamic.BuildLockfile(plan)
				if err := writeJSONFile(flags.LockfilePath, lock); err != nil {
					return fmt.Errorf("write lockfile: %w", err)
				}
			}
			return nil
		},
	}
	bindDynamicFlags(cmd, flags)
	cmd.Flags().BoolVar(&flags.DryRun, "dry-run", false, "Show execution result without mutating machine")
	cmd.Flags().BoolVar(&flags.Resume, "resume", false, "Resume from previous checkpoint state")
	cmd.Flags().StringVar(&flags.StatePath, "state", ".prepare.state.json", "Checkpoint state file path")
	cmd.Flags().StringVar(&flags.LockfilePath, "lockfile", "prepare.lock.json", "Write lockfile after successful run")
	return cmd
}

func newLintCmd() *cobra.Command {
	flags := &dynamicFlags{}
	cmd := &cobra.Command{
		Use:   "lint",
		Short: "Validate manifest syntax and semantic constraints",
		RunE: func(cmd *cobra.Command, args []string) error {
			if flags.ManifestPath == "" {
				flags.ManifestPath = "prepare.yaml"
			}
			manifest, err := dynamic.LoadManifest(flags.ManifestPath)
			if err != nil {
				return err
			}
			catalog := dynamic.BuiltinCatalog()
			profiles := dynamic.BuiltinProfiles()
			if err := dynamic.ValidateManifest(manifest, catalog, profiles); err != nil {
				if flags.OutputJSON {
					_ = dynamic.PrintJSON(map[string]any{"valid": false, "error": err.Error()})
				}
				return err
			}
			if flags.OutputJSON {
				return dynamic.PrintJSON(map[string]any{"valid": true})
			}
			fmt.Println("manifest is valid")
			return nil
		},
	}
	bindDynamicFlags(cmd, flags)
	return cmd
}

func newLockCmd() *cobra.Command {
	flags := &dynamicFlags{}
	cmd := &cobra.Command{
		Use:   "lock",
		Short: "Generate a lockfile from the resolved execution plan",
		RunE: func(cmd *cobra.Command, args []string) error {
			plan, _, err := buildPlan(flags)
			if err != nil {
				return err
			}
			if flags.LockfilePath == "" {
				flags.LockfilePath = "prepare.lock.json"
			}
			if err := writeJSONFile(flags.LockfilePath, dynamic.BuildLockfile(plan)); err != nil {
				return err
			}
			if flags.OutputJSON {
				return dynamic.PrintJSON(map[string]any{"lockfile": flags.LockfilePath})
			}
			fmt.Printf("lockfile generated at %s\n", flags.LockfilePath)
			return nil
		},
	}
	bindDynamicFlags(cmd, flags)
	cmd.Flags().StringVar(&flags.LockfilePath, "lockfile", "prepare.lock.json", "Lockfile path")
	return cmd
}

func bindDynamicFlags(cmd *cobra.Command, flags *dynamicFlags) {
	cmd.Flags().StringVarP(&flags.ManifestPath, "file", "f", "", "Manifest file path (prepare.yaml|prepare.json)")
	cmd.Flags().StringVarP(&flags.Profile, "profile", "p", "", "Profile name to execute")
	cmd.Flags().BoolVar(&flags.OutputJSON, "json", false, "Emit JSON output")
}

func buildPlan(flags *dynamicFlags) (dynamic.Plan, dynamic.Manifest, error) {
	manifest := dynamic.Manifest{APIVersion: "v1", Profiles: map[string]dynamic.Profile{}}
	if flags.ManifestPath != "" {
		path, err := dynamic.DiscoverManifestPath(flags.ManifestPath)
		if err == nil {
			manifest, err = dynamic.LoadManifest(path)
			if err != nil {
				return dynamic.Plan{}, dynamic.Manifest{}, err
			}
		} else if !errors.Is(err, os.ErrNotExist) {
			return dynamic.Plan{}, dynamic.Manifest{}, err
		}
	}

	catalog := dynamic.BuiltinCatalog()
	builtinProfiles := dynamic.BuiltinProfiles()
	if err := dynamic.ValidateManifest(manifest, catalog, builtinProfiles); err != nil {
		return dynamic.Plan{}, dynamic.Manifest{}, err
	}
	tools, err := dynamic.ResolveTools(manifest, flags.Profile, catalog, builtinProfiles)
	if err != nil {
		return dynamic.Plan{}, dynamic.Manifest{}, err
	}
	plan, err := dynamic.BuildPlan(tools, catalog)
	if err != nil {
		return dynamic.Plan{}, dynamic.Manifest{}, err
	}
	return plan, manifest, nil
}

func writeJSONFile(path string, v any) error {
	b, err := jsonMarshalIndent(v)
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}

func jsonMarshalIndent(v any) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number of Go Env Prepare",
		Long:  `All software has versions. This is Go Env Prepare's`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Go Env Prepare v0.0.1")
		},
	}
}

func startPrompt() {
	var selectedOptions []string
	prompt := &survey.MultiSelect{
		Message: "Select some tools:",
		Options: install.GetTitles(),
		Description: func(value string, index int) string {
			return install.GetDescriptionByTitle(value)
		},
		Default: install.GetAlreadyInstalledTools(),
	}
	survey.AskOne(prompt, &selectedOptions)
	for _, option := range selectedOptions {
		installer := install.GetInstallerByTitle(option)
		if installer != nil {
			installer.Install()
		}
	}
	install.FinishAllInstallations()
	shouldRestart := false
	promptRestart := &survey.Confirm{
		Message: "Do you want to restart?",
		Default: false,
	}
	survey.AskOne(promptRestart, &shouldRestart)
	if shouldRestart {
		startPrompt()
		return
	}
	fmt.Println("Bye! 👋")
}
