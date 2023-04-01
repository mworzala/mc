package cli

import "github.com/spf13/cobra"

type Error struct {
	Message string
	Hint    string
	Cmd     *cobra.Command
}

func (e *Error) Error() string {
	return e.Message
}
