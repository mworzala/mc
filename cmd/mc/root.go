package main

import (
	"fmt"
	"os"

	"github.com/mworzala/mc-cli/cmd/mc/account"
	"github.com/mworzala/mc-cli/cmd/mc/java"
	"github.com/mworzala/mc-cli/internal/pkg/cli"
	"github.com/mworzala/mc-cli/internal/pkg/cli/output"
	"github.com/spf13/cobra"
)

var longDescription = `Blah blah blah need to write a longer description
it can have newlines as well :O`

func newRootCmd(app *cli.App) *cobra.Command {
	var outputFormat string

	cmd := &cobra.Command{
		Use:     "mc",
		Version: "0.1.1",
		Short:   "mc is a Minecraft installer and launcher",
		Long:    longDescription,
		PersistentPreRunE: func(_ *cobra.Command, _ []string) (err error) {
			app.Output, err = output.ParseFormat(outputFormat)
			return
		},
	}

	cmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "", "json|yaml|template")

	cmd.AddCommand(account.NewAccountCmd(app))
	cmd.AddCommand(java.NewJavaCmd(app))
	cmd.AddCommand(launchCmd)
	cmd.AddCommand(newInstallCmd(app))

	return cmd
}

func main() {
	app := cli.NewApp()
	rootCmd := newRootCmd(app)

	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
