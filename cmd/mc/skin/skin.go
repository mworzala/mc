package skin

import (
	"github.com/mworzala/mc/internal/pkg/cli"
	"github.com/spf13/cobra"
)

func NewSkinCmd(app *cli.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "skin",
		Short: "Manage Minecraft Skins and Capes",
	}

	cmd.AddCommand(newListCmd(app))
	cmd.AddCommand(newAddCmd(app))
	cmd.AddCommand(newApplyCmd(app))

	return cmd
}
