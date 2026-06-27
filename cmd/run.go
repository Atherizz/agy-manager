package cmd

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/Atherizz/agy-manager/internal/profile"
)

var runCmd = &cobra.Command{
	Use:                "run <profile> -- <cmd>",
	Short:              "Run command with specific profile",
	Args:               cobra.MinimumNArgs(1),
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		prof := args[0]
		var cmdArgs []string
		for i, a := range args {
			if a == "--" {
				cmdArgs = args[i+1:]
				break
			}
		}
		if len(cmdArgs) == 0 {
			return nil
		}
		root, _ := profile.DefaultGeminiRoot()
		p := profile.NewPathResolver(root)
		state, _ := profile.LoadState(p)
		orig := state.ActiveProfile
		if orig != "" {
			profile.StashCurrentProfile(p, orig)
		}
		profile.SwapToProfile(p, prof)
		c := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		c.Stdin, c.Stdout, c.Stderr = os.Stdin, os.Stdout, os.Stderr
		runErr := c.Run()
		profile.StashCurrentProfile(p, prof)
		if orig != "" {
			profile.SwapToProfile(p, orig)
		}
		return runErr
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
