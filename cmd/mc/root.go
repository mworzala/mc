package mc

import (
	"fmt"

	"github.com/mworzala/mc-cli/cmd/mc/account"
	"github.com/mworzala/mc-cli/cmd/mc/java"
	"github.com/mworzala/mc-cli/internal/pkg/cli"
	"github.com/mworzala/mc-cli/internal/pkg/cli/output"
	"github.com/spf13/cobra"
)

var longDescription = `Blah blah blah need to write a longer description
it can have newlines as well :O`

func NewRootCmd(app *cli.App) *cobra.Command {
	var outputFormat string

	versionStr := "dev"
	if app.Build.Commit != "none" {
		versionStr = fmt.Sprintf("%s+%s", app.Build.Version, app.Build.Commit[0:6])
	}

	cmd := &cobra.Command{
		Use:     "mc",
		Version: versionStr,
		Short:   "mc is a Minecraft installer and launcher",
		Long:    longDescription,
		PersistentPreRunE: func(_ *cobra.Command, _ []string) (err error) {
			app.Output, err = output.ParseFormat(outputFormat)
			return
		},
	}

	cmd.SetVersionTemplate(`{{printf "%s" .Version}}`)

	cmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "", "json|yaml|template")

	cmd.AddCommand(account.NewAccountCmd(app))
	cmd.AddCommand(java.NewJavaCmd(app))
	cmd.AddCommand(newLaunchCmd(app))
	cmd.AddCommand(newInstallCmd(app))
	cmd.AddCommand(newVersionCmd(app))

	return cmd
}
