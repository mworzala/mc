package mc

import (
	"encoding/json"

	"github.com/mworzala/mc-cli/internal/pkg/cli"
	"github.com/spf13/cobra"
)

type debugOpts struct {
	app *cli.App
}

func newDebugCmd(app *cli.App) *cobra.Command {
	var o debugOpts

	cmd := &cobra.Command{
		Use:    "debug",
		Hidden: true,
		Args:   cobra.NoArgs,
		RunE: func(_ *cobra.Command, args []string) error {
			o.app = app
			return o.debug(args)
		},
	}

	return cmd
}

func (o *debugOpts) debug(args []string) error {

	res, err := json.MarshalIndent(o.app.Config, "", "  ")
	if err != nil {
		return err
	}
	println(string(res))

	return nil
}
