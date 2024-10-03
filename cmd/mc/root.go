package mc

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/mworzala/mc/cmd/mc/modrinth"

	"github.com/mworzala/mc/cmd/mc/profile"

	"github.com/mworzala/mc/cmd/mc/account"
	"github.com/mworzala/mc/cmd/mc/java"
	"github.com/mworzala/mc/internal/pkg/cli"
	"github.com/mworzala/mc/internal/pkg/cli/output"
	"github.com/spf13/cobra"
)

func NewRootCmd(app *cli.App) *cobra.Command {
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "mc",
		Short: "Minecraft CLI",
		Long:  "Install and manage multiple Minecraft installations from the command line.",
		Example: heredoc.Doc(`
			$ mc account login
			$ mc install 1.20.1 --fabric
			$ mc run 1.20.1-fabric`),
		PersistentPreRunE: func(_ *cobra.Command, _ []string) (err error) {
			app.Output, err = output.ParseFormat(outputFormat)
			return
		},
	}

	cmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "", "json|yaml|template")
	cmd.PersistentFlags().BoolVarP(&app.Config.NonInteractive, "non-interactive", "", false, "disable interactive prompts")

	cmd.AddCommand(account.NewAccountCmd(app))
	cmd.AddCommand(java.NewJavaCmd(app))
	cmd.AddCommand(profile.NewProfileCmd(app))
	cmd.AddCommand(newLaunchCmd(app))
	cmd.AddCommand(newInstallCmd(app))
	cmd.AddCommand(modrinth.NewModrinthCmd(app))
	cmd.AddCommand(newVersionCmd(app))
	cmd.AddCommand(newDebugCmd(app))

	return cmd
}
