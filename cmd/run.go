package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/vishvananda/netlink"
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
	Short: "Run a container with the specified root filesystem and command",
	Long: `Runs a new container with an isolated root filesystem, 
           process ID namespace, and networking. Supports executing specific commands.`,
	Args: cobra.MinimumNArgs(1),
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

func init() {
	RootCmd.AddCommand(RunCmd)
}

// runContainer executes an isolated container environment
func runContainer(rootfs, command string, commandArgs []string) {
	fmt.Println("Running container with root filesystem:", rootfs)

	containerID := generateID()

	// Command to run inside the container
	cmd := exec.Command(command, commandArgs...)

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWNET,
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Setup filesystem before chroot
	setupFilesystem()

	// Change root to the new filesystem
	if err := syscall.Chroot(rootfs); err != nil {
		fmt.Println("Error in chroot:", err)
		return
	}
	if err := os.Chdir("/"); err != nil {
		fmt.Println("Error changing directory:", err)
		return
	}

	// Setup cgroups
	setupCgroups(containerID)

	// Setup networking
	setupNetworking(containerID)

	// Run the container process
	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting container:", err)
		return
	}

	// Store metadata after successful start
	container := Container{ID: containerID, PID: cmd.Process.Pid, Rootfs: rootfs}
	saveContainerMetadata(container)

	// Wait for the process to finish
	if err := cmd.Wait(); err != nil {
		fmt.Println("Container process exited with error:", err)
	}
}

// setupFilesystem mounts /proc and /dev for the container
func setupFilesystem() {
	fmt.Println("Setting up filesystem inside container...")
	syscall.Mount("proc", "/proc", "proc", 0, "")
	os.MkdirAll("/dev", 0755)
	syscall.Mount("tmpfs", "/dev", "tmpfs", 0, "")
}

// setupCgroups configures resource limits for the container
func setupCgroups(containerID string) {
	cgroupPath := filepath.Join("/sys/fs/cgroup/memory", "containGo-"+containerID)
	os.Mkdir(cgroupPath, 0755)
	os.WriteFile(filepath.Join(cgroupPath, "memory.limit_in_bytes"), []byte("268435456"), 0700)
}

// setupNetworking creates a virtual Ethernet pair for networking
func setupNetworking(containerID string) error {
	hostVeth := fmt.Sprintf("veth%s", containerID[:5])
	contVeth := "eth0"

	link := &netlink.Veth{
		LinkAttrs: netlink.LinkAttrs{Name: hostVeth},
		PeerName:  contVeth,
	}
	return netlink.LinkAdd(link)
}

// saveContainerMetadata stores the container metadata to a JSON file
func saveContainerMetadata(container Container) {
	var containers []Container

	// Check if metadata file exists
	if _, err := os.Stat(metadataFile); err == nil {
		file, err := os.Open(metadataFile)
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
	file, err := os.Create(metadataFile)
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
