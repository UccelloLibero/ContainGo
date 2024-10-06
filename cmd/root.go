package cmd

import (
    "github.com/spf13/cobra"
    "fmt"
    "os"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
    Use:   "contain-go",  
    Short: "A lightweight containerization tool built with Go",
    Long: `ContainGo is a simple, lightweight containerization tool built using Go. 
It allows you to run and manage isolated containers using Linux namespaces and cgroups.`,
    // Default action if no subcommand is provided
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Welcome to ContainGo! Use `contain-go --help` to see available commands.")
    },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.go.
func Execute() error {
    return RootCmd.Execute()
}

func init() {
    RootCmd.AddCommand(RunCmd)
    RootCmd.AddCommand(StopCmd)

	// Add a global flag to the root command
    RootCmd.PersistentFlags().String("config", "", "Config file (default is $HOME/.contain-go.yaml)")
}