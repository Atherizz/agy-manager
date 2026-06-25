package cmd

import (

    "fmt"

    "github.com/spf13/cobra"

    "github.com/Atherizz/agy-manager/internal/profile"

)

var listCmd = &cobra.Command{

    Use:     "list",

    Short:   "List all available profiles",

    Aliases: []string{"ls"},

    RunE: func(cmd *cobra.Command, args []string) error {

        root, err := profile.DefaultGeminiRoot()

        if err != nil { return err }

        p := profile.NewPathResolver(root)

        profiles, err := profile.List(p)

        if err != nil { return err }

        if len(profiles) == 0 {

            fmt.Printf("%s No profiles found.\n", infoStyle.Render("ℹ"))

            return nil

        }

        state, err := profile.LoadState(p)

        if err != nil { return err }

        fmt.Printf("%s\n\n", labelStyle.Render("Profiles:"))

        for _, name := range profiles {

            if name == state.ActiveProfile {

                fmt.Printf("  %s %s %s\n", successStyle.Render("▶"), infoStyle.Render(name), labelStyle.Render("(active)"))

            } else {

                fmt.Printf("    %s\n", name)

            }

        }

        return nil

    },

}

func init() { rootCmd.AddCommand(listCmd) }
