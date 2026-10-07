package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

func RegisterWatch(root *cobra.Command) {
	var runID string
	var intervalStr string
	var maxWaitStr string

	watchCmd := &cobra.Command{
		Use:   "watch",
		Short: "Watch a workflow run until completion",
		Run: func(cmd *cobra.Command, args []string) {
			interval := 5 * time.Second
			if intervalStr != "" {
				d, err := strconv.Atoi(intervalStr)
				if err != nil || d < 1 {
					fmt.Fprintf(os.Stderr, "invalid interval: %v\n", err)
					os.Exit(1)
				}
				interval = time.Duration(d) * time.Second
			}

			maxWait := 30 * time.Minute
			if maxWaitStr != "" {
				m, err := strconv.Atoi(maxWaitStr)
				if err != nil || m < 1 {
					fmt.Fprintf(os.Stderr, "invalid max-wait: %v\n", err)
					os.Exit(1)
				}
				maxWait = time.Duration(m) * time.Minute
			}

			deadline := time.Now().Add(maxWait)
			lastStatus := ""

			for time.Now().Before(deadline) {
				out, err := exec.Command("gh", "run", "view", runID, "--json", "status,conclusion,url,name,headSha").CombinedOutput()
				if err != nil {
					fmt.Fprintf(os.Stderr, "gh error: %v\n%s\n", err, out)
					os.Exit(1)
				}

				var result struct {
					Status     string `json:"status"`
					Conclusion string `json:"conclusion"`
					URL        string `json:"url"`
					Name       string `json:"name"`
					HeadSHA    string `json:"headSha"`
				}

				if err := parseJSON(out, &result); err != nil {
					fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
					os.Exit(1)
				}

				statusLine := fmt.Sprintf("[%s] %s (sha: %s)", statusLabel(result.Status, result.Conclusion), result.Name, result.HeadSHA)
				if statusLine != lastStatus {
					fmt.Println(statusLine)
					lastStatus = statusLine
				}

				if result.Status == "completed" {
					fmt.Printf("url: %s\n", result.URL)
					os.Exit(0)
				}

				time.Sleep(interval)
			}

			fmt.Fprintln(os.Stderr, "timeout waiting for workflow")
			os.Exit(1)
		},
	}

	watchCmd.Flags().StringVar(&runID, "run-id", "", "Workflow run ID (required)")
	watchCmd.Flags().StringVar(&intervalStr, "interval", "5", "Polling interval in seconds")
	watchCmd.Flags().StringVar(&maxWaitStr, "max-wait", "30", "Max wait in minutes")
	watchCmd.MarkFlagRequired("run-id")

	root.AddCommand(watchCmd)
}

func statusLabel(status, conclusion string) string {
	switch status {
	case "completed":
		switch conclusion {
		case "success":
			return "Passed"
		case "failure":
			return "Failed"
		default:
			return "Completed"
		}
	case "in_progress":
		return "Running"
	default:
		return status
	}
}

func parseJSON(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
