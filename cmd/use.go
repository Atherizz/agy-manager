package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Atherizz/agy-manager/internal/profile"
)

var useCmd = &cobra.Command{
	Use:   "use [profile-name]",
	Short: "Switch to a different profile",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		root, _ := profile.DefaultGeminiRoot()
		p := profile.NewPathResolver(root)

		if len(args) == 0 {
			return runInteractiveSelector(p)
		}

		name := args[0]
		if !profile.Exists(p, name) {
			return fmt.Errorf("%s Profile not found", errorStyle.Render("✗"))
		}
		state, _ := profile.LoadState(p)
		if state.ActiveProfile == name {
			return nil
		}
		if state.ActiveProfile != "" && profile.Exists(p, state.ActiveProfile) {
			fmt.Printf("  %s Stashing %s...\n", labelStyle.Render("→"), state.ActiveProfile)
			profile.StashCurrentProfile(p, state.ActiveProfile)
		}
		fmt.Printf("  %s Loading %s...\n", labelStyle.Render("→"), name)
		if err := profile.SwapToProfile(p, name); err != nil {
			return err
		}
		state.ActiveProfile = name
		profile.SaveState(p, state)
		fmt.Printf("\n%s Switched to profile %s\n", successStyle.Render("✓"), infoStyle.Render(name))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(useCmd)
}