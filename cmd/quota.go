package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/Atherizz/agy-manager/internal/profile"
)

var quotaCmd = &cobra.Command{
	Use:   "quota [profile-name]",
	Short: "Check AI model usage quota for a profile (or all profiles)",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		root, _ := profile.DefaultGeminiRoot()
		p := profile.NewPathResolver(root)
		state, _ := profile.LoadState(p)

		var profilesToCheck []string
		if len(args) == 1 {
			if !profile.Exists(p, args[0]) {
				return fmt.Errorf("%s Profile %q not found", errorStyle.Render("✗"), args[0])
			}
			profilesToCheck = []string{args[0]}
		} else {
			profilesToCheck, _ = profile.List(p)
		}

		if len(profilesToCheck) == 0 {
			fmt.Println("No profiles found.")
			return nil
		}

		for i, name := range profilesToCheck {
			fmt.Printf("  %s Fetching quota for %s...\n", labelStyle.Render("→"), infoStyle.Render(name))

			quotas, err := profile.GetQuota(name, state.ActiveProfile)
			if err != nil {
				fmt.Printf("    %s Failed: %v\n\n", errorStyle.Render("✗"), err)
				continue
			}

			if len(quotas) == 0 {
				fmt.Printf("    %s No supported models found or quota data unavailable.\n\n", errorStyle.Render("✗"))
				continue
			}

			var models []string
			for k := range quotas {
				models = append(models, k)
			}
			sort.Strings(models)

			for _, model := range models {
				percentage := quotas[model] * 100
				displayName := strings.ReplaceAll(model, "models/", "")
				
				pctStr := fmt.Sprintf("%.0f%%", percentage)
				if percentage < 10 {
					pctStr = errorStyle.Render(pctStr)
				} else if percentage < 40 {
					pctStr = labelStyle.Render(pctStr)
				} else {
					pctStr = successStyle.Render(pctStr)
				}

				fmt.Printf("    %-35s %s\n", displayName, pctStr)
			}
			
			// Add a blank line between profiles unless it's the last one
			if i < len(profilesToCheck)-1 {
				fmt.Println()
			}
		}

		fmt.Println()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(quotaCmd)
}
