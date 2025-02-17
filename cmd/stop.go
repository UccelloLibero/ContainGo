package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/spf13/cobra"
)

// Metadata file path
const metadataFile = "/var/lib/contain-go/containers.json"

// StopCmd represents the stop command
var StopCmd = &cobra.Command{
	Use:   "stop [pid]",
	Short: "Stop a running container",
	Long:  `Stops a running container by sending SIGKILL and removing its resources.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pid, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Invalid PID. Please provide a valid process ID.")
			return
		}
		stopContainer(pid)
	},
}

func init() {
	RootCmd.AddCommand(StopCmd)
}

func stopContainer(pid int) {
	fmt.Printf("Stopping container with PID: %d\n", pid)

	proc, err := os.FindProcess(pid)
	if err != nil {
		fmt.Println("Error finding process:", err)
		return
	}

	if err := proc.Signal(syscall.SIGKILL); err != nil {
		fmt.Println("Error sending SIGKILL:", err)
	} else {
		fmt.Println("Container stopped.")
	}

	// Find the container ID associated with the PID
	containerID := retrieveContainerID(pid)

	if containerID != "" {
		fmt.Println("Removing container resources for:", containerID)

		// Clean up cgroups
		cgroupPath := filepath.Join("/sys/fs/cgroup/memory", "containGo-"+containerID)
		os.RemoveAll(cgroupPath)

		// Remove container metadata
		removeContainerMetadata(containerID)
	}
}

// retrieveContainerID finds the container ID based on the given PID
func retrieveContainerID(pid int) string {
	file, err := os.Open(metadataFile)
	if err != nil {
		fmt.Println("Error opening metadata file:", err)
		return ""
	}
	defer file.Close()

	var containers []Container

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&containers); err != nil {
		fmt.Println("Error decoding container metadata:", err)
		return ""
	}

	for _, container := range containers {
		if container.PID == pid {
			return container.ID
		}
	}

	fmt.Println("No container found for PID:", pid)
	return ""
}

// removeContainerMetadata deletes a container's metadata
func removeContainerMetadata(containerID string) {
	file, err := os.Open(metadataFile)
	if err != nil {
		fmt.Println("Error opening metadata file:", err)
		return
	}
	defer file.Close()

	var containers []Container
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&containers); err != nil {
		fmt.Println("Error decoding container metadata:", err)
		return
	}

	// Remove the container from the list
	var updatedContainers []Container
	for _, container := range containers {
		if container.ID != containerID {
			updatedContainers = append(updatedContainers, container)
		}
	}

	// Write the updated list back to the file
	file, err = os.Create(metadataFile)
	if err != nil {
		fmt.Println("Error creating metadata file:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(updatedContainers); err != nil {
		fmt.Println("Error encoding metadata:", err)
	}
}
