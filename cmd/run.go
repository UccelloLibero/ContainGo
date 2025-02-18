package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"

	"github.com/spf13/cobra"
)

// Container struct to store metadata
type Container struct {
	ID     string `json:"id"`
	PID    int    `json:"pid"`
	Rootfs string `json:"rootfs"`
}

// generateID creates a unique container ID
func generateID() string {
	return fmt.Sprintf("%x", time.Now().UnixNano())
}

// RunCmd represents the run command
var RunCmd = &cobra.Command{
	Use:   "run [rootfs] [command] [args...]",
	Short: "Run a process inside an isolated root filesystem",
	Long:  "Runs a command inside a new root filesystem using chroot. Works on macOS and Linux.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		rootfs := args[0]

		// Default to a shell if no command is provided
		command := "/bin/sh"
		commandArgs := []string{}
		if len(args) > 1 {
			command = args[1]
			commandArgs = args[2:]
		}

		runContainer(rootfs, command, commandArgs)
	},
}

// runContainer runs a process inside an isolated root filesystem
func runContainer(rootfs, command string, commandArgs []string) {
	fmt.Println("Running container with root filesystem:", rootfs)

	containerID := generateID()

	// Fork and run the container process
	cmd := exec.Command(command, commandArgs...)

	// Skip namespace isolation (not supported on macOS/Windows)
	fmt.Println("No namespace isolation (macOS & Linux-compatible)")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Perform chroot (only for Linux/macOS)
	if runtime.GOOS != "windows" {
		if err := syscall.Chroot(rootfs); err != nil {
			fmt.Println("Error in chroot:", err)
			return
		}
		if err := os.Chdir("/"); err != nil {
			fmt.Println("Error changing directory:", err)
			return
		}

		// Setup filesystem only on Linux
		if runtime.GOOS == "linux" {
			setupFilesystem()
		} else {
			fmt.Println("Skipping filesystem setup. macOS does not support /proc mounting.")
		}
	} else {
		fmt.Println("Skipping chroot: Windows does not support this operation.")
	}

	// Run the container process
	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting container:", err)
		return
	}

	// Store metadata
	container := Container{ID: containerID, PID: cmd.Process.Pid, Rootfs: rootfs}
	saveContainerMetadata(container)

	// Wait for the process to finish
	if err := cmd.Wait(); err != nil {
		fmt.Println("Container process exited with error:", err)
	}
}

// setupFilesystem mounts /proc (only on Linux)
func setupFilesystem() {
	if runtime.GOOS != "linux" {
		return
	}
	fmt.Println("Setting up filesystem inside container...")

	// Mount /proc
	if err := unix.Mount("proc", "/proc", 0, unsafe.Pointer(nil)); err != nil {
		fmt.Println("Error mounting /proc:", err)
	}
}

// saveContainerMetadata stores the container metadata
func saveContainerMetadata(container Container) {
	var containers []Container

	// Check if metadata file exists
	if _, err := os.Stat("metadata.json"); err == nil {
		file, err := os.Open("metadata.json")
		if err != nil {
			fmt.Println("Error opening metadata file:", err)
			return
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&containers); err != nil {
			fmt.Println("Error decoding existing metadata:", err)
			return
		}
	}

	// Append the new container
	containers = append(containers, container)

	// Write updated metadata back to file
	file, err := os.Create("metadata.json")
	if err != nil {
		fmt.Println("Error creating metadata file:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(containers); err != nil {
		fmt.Println("Error encoding metadata:", err)
	}
}
