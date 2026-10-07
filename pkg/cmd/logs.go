package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func RegisterLogs(root *cobra.Command) {
	var runID string
	var errorsOnly bool
	var limitStr string

	logsCmd := &cobra.Command{
		Use:   "logs",
		Short: "Fetch workflow run logs",
		Run: func(cmd *cobra.Command, args []string) {
			ghArgs := []string{"run", "view", runID, "--log-failed"}
			out, err := exec.Command("gh", ghArgs...).CombinedOutput()
			if err != nil {
				fmt.Fprintf(os.Stderr, "gh error: %v\n%s\n", err, out)
				os.Exit(1)
			}

			logs := string(out)
			if errorsOnly {
				logs = extractErrors(logs)
			}

			if limitStr != "" {
				limit := 0
				if _, err := fmt.Sscanf(limitStr, "%d", &limit); err == nil && limit > 0 {
					logs = limitLines(logs, limit)
				}
			}

			fmt.Print(logs)
		},
	}

	logsCmd.Flags().StringVar(&runID, "run-id", "", "Workflow run ID (required)")
	logsCmd.Flags().BoolVar(&errorsOnly, "errors-only", false, "Show only error lines")
	logsCmd.Flags().StringVar(&limitStr, "limit", "", "Max lines to return")
	logsCmd.MarkFlagRequired("run-id")

	root.AddCommand(logsCmd)
}

func extractErrors(logs string) string {
	var result []string
	for _, line := range strings.Split(logs, "\n") {
		lower := strings.ToLower(line)
		if strings.Contains(lower, "error") ||
			strings.Contains(lower, "fail") ||
			strings.Contains(lower, "fatal") ||
			strings.Contains(lower, "panic") ||
			strings.Contains(lower, "exception") ||
			strings.Contains(lower, "traceback") ||
			strings.Contains(lower, "exit code") {
			result = append(result, line)
		}
	}
	if len(result) == 0 {
		return "No errors detected."
	}
	return strings.Join(result, "\n")
}

func limitLines(text string, max int) string {
	lines := strings.Split(text, "\n")
	if len(lines) <= max {
		return text
	}
	return strings.Join(lines[len(lines)-max:], "\n")
}
