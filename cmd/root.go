package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without subcommands
var RootCmd = &cobra.Command{
	Use:   "contain-go",
	Short: "A lightweight containerization tool built with Go",
	Long: `ContainGo allows you to run and manage isolated containers.
It provides filesystem isolation using chroot and executes commands within a containerized environment.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to ContainGo! Use `contain-go --help` to see available commands.")
	},
}

// Execute runs the root command
func Execute() error {
	return RootCmd.Execute()
}

func init() {
	// Ensure these commands are always added
	RootCmd.AddCommand(RunCmd)
	RootCmd.AddCommand(StopCmd)
	// Check if the ListCmd is added to the RootCmd
	if !RootCmd.HasSubCommands() {
		RootCmd.AddCommand(ListCmd)
	}

	RootCmd.PersistentFlags().String("config", "", "Config file (default is $HOME/.contain-go.yaml)")
}
