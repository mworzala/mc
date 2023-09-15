package profile

import (
	"github.com/mworzala/mc-cli/internal/pkg/cli"
	"github.com/spf13/cobra"
)

func NewProfileCmd(app *cli.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profile",
		Short: "Manage profiles (Minecraft installations)",
	}

	cmd.AddCommand(newListCmd(app))

	return cmd
}
