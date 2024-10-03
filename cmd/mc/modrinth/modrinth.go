package modrinth

import (
	"github.com/mworzala/mc/internal/pkg/cli"
	"github.com/spf13/cobra"
)

func NewModrinthCmd(app *cli.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "modrinth",
		Short: "Query the modrinth API directly",
	}

	cmd.AddCommand(newSearchCmd(app))

	return cmd
}
