package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without subcommands
var RootCmd = &cobra.Command{
	Use:   "contain-go",
	Short: "A lightweight containerization tool built with Go",
	Long: `ContainGo is a simple, lightweight containerization tool. 
It allows you to run and manage isolated containers using Linux namespaces and cgroups.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to ContainGo! Use `contain-go --help` to see available commands.")
	},
}

// Execute adds all child commands
func Execute() error {
	return RootCmd.Execute()
}

func init() {
	RootCmd.AddCommand(RunCmd)
	RootCmd.AddCommand(StopCmd)
	RootCmd.PersistentFlags().String("config", "", "Config file (default is $HOME/.contain-go.yaml)")
}
