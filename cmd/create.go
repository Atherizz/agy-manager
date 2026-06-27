package cmd

import (
	"fmt"
	"regexp"

	"github.com/spf13/cobra"

	"github.com/Atherizz/agy-manager/internal/profile"
)

var validName = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]*$`)

var createCmd = &cobra.Command{
	Use:   "create <profile-name>",
	Short: "Create a new Antigravity CLI profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		if !validName.MatchString(name) {
			return fmt.Errorf("%s Invalid name %q", errorStyle.Render("✗"), name)
		}
		root, err := profile.DefaultGeminiRoot()
		if err != nil {
			return err
		}
		p := profile.NewPathResolver(root)
		if profile.Exists(p, name) {
			return fmt.Errorf("%s Profile %q already exists", errorStyle.Render("✗"), name)
		}
		if err := profile.Create(p, name); err != nil {
			return err
		}
		fmt.Printf("%s Created profile %s\n  %s %s\n",
			successStyle.Render("✓"), infoStyle.Render(name),
			labelStyle.Render("Location:"), p.ProfileDir(name))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}