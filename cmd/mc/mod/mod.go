package mod

import (
	"errors"

	"github.com/mworzala/mc/internal/pkg/cli"
	"github.com/spf13/cobra"
)

var (
	errNoProfile = errors.New("no profile provided")
	errNoModId   = errors.New("no mod id provided")
)

func NewModProfileCmd(app *cli.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mod",
		Short: "Manage mods for a profile",
	}

	cmd.AddCommand(newModInstallCmd(app))

	return cmd
}
