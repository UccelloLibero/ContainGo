package cmd

import (
    "fmt"
    "os"
    "strconv"
    "syscall"

    "github.com/spf13/cobra"
)

// StopCmd represents the stop command
var StopCmd = &cobra.Command{
    Use:   "stop [pid]",  // Command format expects a process ID
    Short: "Stop a running container",
    Long: `This command stops a running container by sending a SIGKILL signal 
to the container's process. You need to provide the process ID (PID) of the 
container to stop it.`,
    Args: cobra.ExactArgs(1),  // Ensure exactly one argument is provided (the PID)
    Run: func(cmd *cobra.Command, args []string) {
        pid, err := strconv.Atoi(args[0])  // Convert the PID from string to integer
        if err != nil {
            fmt.Println("Invalid PID. Please provide a valid process ID.")
            return
        }
        stopContainer(pid)
    },
}

func init() {
    // Add the stop command to the root command
    RootCmd.AddCommand(StopCmd)
}

// stopContainer sends a SIGKILL to the process with the given PID
func stopContainer(pid int) {
    fmt.Printf("Stopping container with PID: %d\n", pid)

    // Find the process by PID
    proc, err := os.FindProcess(pid)
    if err != nil {
        fmt.Println("Error finding process:", err)
        return
    }

    // Send the SIGKILL signal to terminate the process
    if err := proc.Signal(syscall.SIGKILL); err != nil {
        fmt.Println("Error sending SIGKILL signal:", err)
    } else {
        fmt.Println("Container stopped.")
    }
}