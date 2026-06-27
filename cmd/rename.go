package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/Atherizz/agy-manager/internal/profile"
)

var renameCmd = &cobra.Command{
	Use:   "rename <old-name> <new-name>",
	Short: "Rename a profile",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		oldName, newName := args[0], args[1]
		root, _ := profile.DefaultGeminiRoot()
		p := profile.NewPathResolver(root)

		if !profile.Exists(p, oldName) {
			return fmt.Errorf("%s Profile %q not found", errorStyle.Render("✗"), oldName)
		}
		if !validName.MatchString(newName) {
			return fmt.Errorf("%s Invalid name %q", errorStyle.Render("✗"), newName)
		}
		if profile.Exists(p, newName) {
			return fmt.Errorf("%s Profile %q already exists", errorStyle.Render("✗"), newName)
		}

		state, _ := profile.LoadState(p)
		if state.ActiveProfile == oldName {
			// Active profile: vault folder exists but cred is at main slot, not profile slot.
			// Just rename the folder and update state — no WCM rename needed.
			if err := os.Rename(p.ProfileDir(oldName), p.ProfileDir(newName)); err != nil {
				return err
			}
			state.ActiveProfile = newName
			profile.SaveState(p, state)
		} else {
			if err := profile.Rename(p, oldName, newName); err != nil {
				return err
			}
		}

		fmt.Printf("%s Renamed %s → %s\n",
			successStyle.Render("✓"),
			infoStyle.Render(oldName),
			infoStyle.Render(newName))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(renameCmd)
}
