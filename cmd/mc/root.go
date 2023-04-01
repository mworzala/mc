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

	cmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "default", "json|yaml|template")

	cmd.AddCommand(launchCmd)
	cmd.AddCommand(installCmd)
	cmd.AddCommand(account.NewAccountCmd(app))
	cmd.AddCommand(java.Cmd)

	return cmd
}

//var (
//	outputFlag string
//)

//func init() {
//	rootCmd.PersistentFlags().StringVarP(&outputFlag, "output", "o", "", "json|yaml")
//
//	rootCmd.AddCommand(launchCmd)
//	rootCmd.AddCommand(installCmd)
//	rootCmd.AddCommand(account.Cmd)
//	rootCmd.AddCommand(java.Cmd)
//}

func preSetupFlags(_ *cobra.Command, _ []string) error {

	//	outputOverride, err := cmd.PersistentFlags().GetString("output")
	//	if err != nil {
	//		println("ROOT ERRR")
	//		return err
	//	}
	//
	//	outputOverride = strings.ToLower(outputOverride)
	//	if _, ok := appModel.OutputFormatValidationMap[outputOverride]; !ok {
	//		return fmt.Errorf("invalid output format: %s", outputOverride)
	//	}

	return nil
}

func main() {
	app := cli.NewApp()
	rootCmd := newRootCmd(app)

	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
