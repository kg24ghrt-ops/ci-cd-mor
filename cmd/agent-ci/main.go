package main

import (
	"os"

	"github.com/kg24ghrt-ops/ci-cd-mor/pkg/cmd"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "agent-ci",
		Short: "Agent-focused CI/CD CLI for GitHub Actions",
	}

	cmd.RegisterStart(rootCmd)
	cmd.RegisterWatch(rootCmd)
	cmd.RegisterLogs(rootCmd)
	cmd.RegisterStop(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
