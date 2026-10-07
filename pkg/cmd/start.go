package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func RegisterStart(root *cobra.Command) {
	var workflowFile string
	var branch string

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start a GitHub Actions workflow",
		Run: func(cmd *cobra.Command, args []string) {
			ghArgs := []string{"workflow", "run", workflowFile, "--branch", branch, "--json", "runId,status,url,conclusion,headSha,headBranch"}

			out, err := exec.Command("gh", ghArgs...).CombinedOutput()
			if err != nil {
				fmt.Fprintf(os.Stderr, "gh error: %v\n%s\n", err, out)
				os.Exit(1)
			}

			var result struct {
				RunID      string `json:"runId"`
				Status     string `json:"status"`
				URL        string `json:"url"`
				Conclusion string `json:"conclusion"`
				HeadSHA    string `json:"headSha"`
				HeadBranch string `json:"headBranch"`
			}

			if err := json.Unmarshal(out, &result); err != nil {
				fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
				os.Exit(1)
			}

			response := map[string]interface{}{
				"runId":      result.RunID,
				"status":     result.Status,
				"url":        result.URL,
				"conclusion": result.Conclusion,
				"headSha":    result.HeadSHA,
				"headBranch": result.HeadBranch,
			}

			data, _ := json.MarshalIndent(response, "", "  ")
			fmt.Println(string(data))
		},
	}

	startCmd.Flags().StringVar(&workflowFile, "workflow", "", "Workflow file (required)")
	startCmd.Flags().StringVar(&branch, "branch", "main", "Git branch")
	startCmd.MarkFlagRequired("workflow")

	root.AddCommand(startCmd)
}
