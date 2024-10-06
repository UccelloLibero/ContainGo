package cmd

import (
    "fmt"
    "os"
    "os/exec"
    "syscall"
    "github.com/spf13/cobra"
)

// RunCmd represents the run command
var RunCmd = &cobra.Command{
    Use:   "run [rootfs]",  // Command format
    Short: "Run a container with the specified root filesystem",
    Long: `This command runs a new container using the specified root filesystem directory. 
It sets up process and filesystem isolation using Linux namespaces.`,
    Args: cobra.ExactArgs(1),  // Expect exactly 1 argument (the root filesystem path)
    Run: func(cmd *cobra.Command, args []string) {
        rootfs := args[0]
        runContainer(rootfs)
    },
}

func init() {
    // Add the run command to the root command
    RootCmd.AddCommand(RunCmd)
}

// runContainer runs the container with the specified root filesystem
func runContainer(rootfs string) {
    fmt.Println("Running container with root filesystem:", rootfs)

    // Step 1: Command to run inside the container (e.g., /bin/sh)
    cmd := exec.Command("/bin/sh")

    // Step 2: Set up Linux namespaces for process, filesystem, UTS, and networking isolation
    cmd.SysProcAttr = &syscall.SysProcAttr{
        Cloneflags: syscall.CLONE_NEWNS   |  // Filesystem isolation
                    syscall.CLONE_NEWPID  |  // Process ID isolation
                    syscall.CLONE_NEWUTS  |  // Hostname isolation
                    syscall.CLONE_NEWNET,    // Networking isolation
    }

    // Step 3: Hook up standard input, output, and error to the terminal
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    // Step 4: Change root filesystem using chroot and change to the new root directory
    if err := syscall.Chroot(rootfs); err != nil {
        fmt.Println("Error in chroot:", err)
        return
    }
    if err := os.Chdir("/"); err != nil {
        fmt.Println("Error changing directory:", err)
        return
    }

    // Step 5: Run the shell inside the isolated container environment
    if err := cmd.Run(); err != nil {
        fmt.Println("Error running container:", err)
    }
}