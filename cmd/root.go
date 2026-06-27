package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
)

var rootCmd = &cobra.Command{
	Use:   "agym",
	Short: "Antigravity Identity Manager",
	Long: `agym — Manage multiple Antigravity CLI profiles on one machine.

Each profile isolates credentials, settings, and conversation history
while sharing global components like plugins and skills.`,
}

// Execute is the entry point called from main.go.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}