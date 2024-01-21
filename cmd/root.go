package cmd

import (
	"fmt"
	"os"
)

func Execute() {
	commands := NewCommands()
	if err := commands.RootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
