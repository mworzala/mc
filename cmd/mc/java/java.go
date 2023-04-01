package java

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:   "java",
	Short: "Manage java installations",
}

func init() {
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(defaultCmd)
}
