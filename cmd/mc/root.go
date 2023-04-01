package main

import (
	"fmt"
	"os"

	"github.com/mworzala/mc-cli/cmd/mc/account"
	"github.com/mworzala/mc-cli/cmd/mc/java"
	"github.com/spf13/cobra"
)

var longDescription = `Blah blah blah need to write a longer description
it can have newlines as well :O`

var rootCmd = &cobra.Command{
	Use:     "mc",
	Version: "0.1.1",
	Short:   "mc is a Minecraft installer and launcher",
	Long:    longDescription,
}

func init() {
	rootCmd.AddCommand(launchCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(account.Cmd)
	rootCmd.AddCommand(java.Cmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
