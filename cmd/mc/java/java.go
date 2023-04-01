package java

import (
	"github.com/mworzala/mc-cli/internal/pkg/cli"
	"github.com/spf13/cobra"
)

func NewJavaCmd(app *cli.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "java",
		Short: "Manage java installations",
	}

	cmd.AddCommand(newListCmd(app))
	cmd.AddCommand(newDefaultCmd(app))

	return cmd
}
