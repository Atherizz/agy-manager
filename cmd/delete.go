package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Atherizz/agy-manager/internal/profile"
)

var forceDelete bool

var deleteCmd = &cobra.Command{
	Use:   "delete <profile-name>",
	Short: "Delete a profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		root, _ := profile.DefaultGeminiRoot()
		p := profile.NewPathResolver(root)
		if !profile.Exists(p, name) {
			return fmt.Errorf("Profile %s not found", name)
		}
		state, _ := profile.LoadState(p)
		if state.ActiveProfile == name {
			return fmt.Errorf("Cannot delete active profile")
		}
		if !forceDelete {
			fmt.Printf("Delete profile %s? [y/N]: ", name)
			reader := bufio.NewReader(os.Stdin)
			ans, _ := reader.ReadString('\n')
			if strings.TrimSpace(strings.ToLower(ans)) != "y" {
				return nil
			}
		}
		os.RemoveAll(p.ProfileDir(name))
		fmt.Printf("%s Deleted %s\n", successStyle.Render("✓"), name)
		return nil
	},
}

func init() {
	deleteCmd.Flags().BoolVarP(&forceDelete, "force", "f", false, "Skip confirmation")
	rootCmd.AddCommand(deleteCmd)
}
