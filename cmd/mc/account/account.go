package account

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:   "account",
	Short: "Manage accounts or log into a new one",
}

func init() {
	Cmd.AddCommand(loginCmd)
	Cmd.AddCommand(defaultCmd)
}
