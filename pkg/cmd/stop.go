package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func RegisterStop(root *cobra.Command) {
	var runID string

	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Cancel a running workflow",
		Run: func(cmd *cobra.Command, args []string) {
			out, err := exec.Command("gh", "run", "cancel", runID).CombinedOutput()
			if err != nil {
				fmt.Fprintf(os.Stderr, "gh error: %v\n%s\n", err, out)
				os.Exit(1)
			}

			type response struct {
				RunID      string `json:"runId"`
				Status     string `json:"status"`
				Conclusion string `json:"conclusion"`
				URL        string `json:"url"`
				Name       string `json:"name"`
			}

			var result response
			if err := json.Unmarshal(out, &result); err != nil {
				response := map[string]string{
					"runId":  runID,
					"status": "cancelled",
				}
				data, _ := json.MarshalIndent(response, "", "  ")
				fmt.Println(string(data))
				return
			}

			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
		},
	}

	stopCmd.Flags().StringVar(&runID, "run-id", "", "Workflow run ID (required)")
	stopCmd.MarkFlagRequired("run-id")

	root.AddCommand(stopCmd)
}
