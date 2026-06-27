package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Atherizz/agy-manager/internal/profile"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show active profile",
	RunE: func(cmd *cobra.Command, args []string) error {
		root, _ := profile.DefaultGeminiRoot()
		state, _ := profile.LoadState(profile.NewPathResolver(root))
		if state.ActiveProfile == "" {
			fmt.Println("No active profile.")
		} else {
			fmt.Printf("%s Active profile: %s\n", successStyle.Render("▶"), infoStyle.Render(state.ActiveProfile))
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}