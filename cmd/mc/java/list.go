package java

import "github.com/spf13/cobra"

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List discovered java installations",
	Args:  cobra.NoArgs,
	RunE:  handleList,
}

func handleList(_ *cobra.Command, _ []string) error {
	return nil
}
