package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/Atherizz/agy-manager/internal/profile"
)

var infoCmd = &cobra.Command{
	Use:   "info <profile-name>",
	Short: "Show details about a profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		root, _ := profile.DefaultGeminiRoot()
		p := profile.NewPathResolver(root)

		if !profile.Exists(p, name) {
			return fmt.Errorf("%s Profile %q not found", errorStyle.Render("✗"), name)
		}

		state, _ := profile.LoadState(p)
		isActive := state.ActiveProfile == name

		email := profile.GetAccountEmail(name, state.ActiveProfile)
		isEmailLike := len(email) > 0 && email != "antigravity"
		hasCredential := profile.HasWindowsCred(name, state.ActiveProfile)

		fmt.Printf("\n%s\n\n", infoStyle.Render(name))

		status := "inactive"
		if isActive {
			status = successStyle.Render("active")
		}
		fmt.Printf("  %s %s\n", labelStyle.Render("Status:  "), status)

		if isEmailLike {
			fmt.Printf("  %s %s\n", labelStyle.Render("Account: "), email)
		} else if hasCredential {
			fmt.Printf("  %s %s\n", labelStyle.Render("Account: "), labelStyle.Render("credential stored"))
		} else {
			fmt.Printf("  %s %s\n", labelStyle.Render("Account: "), labelStyle.Render("not logged in"))
		}

		fmt.Printf("  %s %s\n", labelStyle.Render("Vault:   "), p.ProfileDir(name))

		// Show which files are in the vault
		var storedFiles []string
		for _, f := range p.IsolatedFiles() {
			path := p.ProfileDir(name) + string(os.PathSeparator) + f
			if _, err := os.Stat(path); err == nil {
				storedFiles = append(storedFiles, f)
			}
		}
		if len(storedFiles) > 0 {
			fmt.Printf("  %s %s\n", labelStyle.Render("Files:   "), storedFiles[0])
			for _, f := range storedFiles[1:] {
				fmt.Printf("             %s\n", f)
			}
		} else if isActive {
			fmt.Printf("  %s %s\n", labelStyle.Render("Files:   "), labelStyle.Render("in root (profile is active)"))
		} else {
			fmt.Printf("  %s %s\n", labelStyle.Render("Files:   "), labelStyle.Render("none"))
		}

		fmt.Println()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
