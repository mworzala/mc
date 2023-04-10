package profile

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/mworzala/mc-cli/internal/pkg/cli"
	"github.com/spf13/cobra"
)

func NewProfileCmd(app *cli.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profile",
		Short: "Install or manage profiles",
		Long: heredoc.Doc(`
			Blah blah blah todo
			AllProfiles are the same thing as instances in multimc/polymc/prism
`),
	}

	cmd.AddCommand(newListCmd(app))

	return cmd
}
