package cmd

import (
	"fmt"

	"github.com/charmbracelet/huh"

	"github.com/Atherizz/agy-manager/internal/profile"
)

func runInteractiveSelector(p *profile.PathResolver) error {
	profiles, err := profile.List(p)
	if err != nil {
		return err
	}
	if len(profiles) == 0 {
		fmt.Printf("%s No profiles found. Run %s first.\n",
			infoStyle.Render("ℹ"),
			infoStyle.Render("agym create <name>"))
		return nil
	}

	state, _ := profile.LoadState(p)

	options := make([]huh.Option[string], len(profiles))
	for i, name := range profiles {
		label := name
		if name == state.ActiveProfile {
			label = name + " " + labelStyle.Render("(active)")
		}
		options[i] = huh.NewOption(label, name)
	}

	var selected string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Switch Profile").
				Options(options...).
				Value(&selected),
		),
	)

	if err := form.Run(); err != nil {
		return nil // user pressed Esc/Ctrl+C
	}

	if selected == state.ActiveProfile {
		fmt.Printf("%s Already on profile %s\n",
			infoStyle.Render("ℹ"), infoStyle.Render(selected))
		return nil
	}

	if state.ActiveProfile != "" && profile.Exists(p, state.ActiveProfile) {
		fmt.Printf("  %s Stashing %s...\n", labelStyle.Render("→"), state.ActiveProfile)
		profile.StashCurrentProfile(p, state.ActiveProfile)
	}

	fmt.Printf("  %s Loading %s...\n", labelStyle.Render("→"), selected)
	if err := profile.SwapToProfile(p, selected); err != nil {
		return err
	}

	state.ActiveProfile = selected
	profile.SaveState(p, state)

	fmt.Printf("\n%s Switched to profile %s\n",
		successStyle.Render("✓"), infoStyle.Render(selected))
	return nil
}
