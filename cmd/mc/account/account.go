package account

import (
	"github.com/mworzala/mc-cli/internal/pkg/cli"
	"github.com/spf13/cobra"
)

func NewAccountCmd(app *cli.App) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "account",
		Short: "Manage accounts or log into a new one",
	}

	cmd.AddCommand(newLoginCmd(app))
	cmd.AddCommand(newDefaultCmd(app))
	cmd.AddCommand(newTokenCmd(app))

	return cmd
}
