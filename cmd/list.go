package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// ListCmd represents the list command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all running containers",
	Long:  `Display all running containers along with their process IDs and PIDs.`,
	Run: func(cmd *cobra.Command, args []string) {
		listContainers()
	},
}

func init() {
	RootCmd.AddCommand(ListCmd)
}

// listContainers reads and prints container metadata from the JSON file
func listContainers() {
	file, err := os.Open(metadataFile)
	if err != nil {
		fmt.Println("No running containers.")
		return
	}
	defer file.Close()

	var containers []Container
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&containers); err != nil {
		fmt.Println("No active containers.")
		return
	}

	fmt.Println("\n📦 Running Containers:")
	for _, container := range containers {
		fmt.Printf("🆔 ID: %s | 🏷️ PID: %d | 📂 RootFS: %s\n", container.ID, container.PID, container.Rootfs)
	}
}
