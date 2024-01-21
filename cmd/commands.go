package cmd

import (
	"felipewom/go-env-prepare/cmd/install"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

type Commands struct {
	RootCmd *cobra.Command
}

func NewCommands() *Commands {
	rootCmd := &Commands{
		RootCmd: newInstallCmd(),
	}
	rootCmd.RootCmd.AddCommand(newVersionCmd())
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
