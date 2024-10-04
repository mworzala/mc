package skin

import (
	"github.com/mworzala/mc/internal/pkg/cli"
	"github.com/spf13/cobra"
)

func NewSkinCmd(app *cli.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "skin",
		Short: "Manage Minecraft skins and capes",
	}

	var account string

	cmd.Flags().StringVar(&account, "account", "", "Account to use")

	cmd.AddCommand(newListCmd(app))
	cmd.AddCommand(newAddCmd(app, account))
	cmd.AddCommand(newApplyCmd(app, account))

	return cmd
}
