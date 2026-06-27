package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of agym",
	Long:  `All software has versions. This is agym's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("agym version %s, commit %s\n", version, commit)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}