package skin

import (
	"github.com/mworzala/mc/internal/pkg/cli"
	"github.com/spf13/cobra"
)

type applySkinOpts struct {
	app *cli.App

	account string
}

func newApplyCmd(app *cli.App) *cobra.Command {
	var o applySkinOpts

	cmd := &cobra.Command{
		Use:     "apply",
		Short:   "Apply a saved skin",
		Aliases: []string{"set"},
		Args: func(cmd *cobra.Command, args []string) error {
			o.app = app
			return o.validateArgs(cmd, args)
		},
		RunE: func(_ *cobra.Command, args []string) error {
			o.app = app
			return o.execute(args)
		},
	}

	cmd.Flags().StringVar(&o.account, "account", "", "Account to use")

	return cmd
}

func (o *applySkinOpts) validateArgs(cmd *cobra.Command, args []string) (err error) {
	if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
		return err
	}

	return nil
}

func (o *applySkinOpts) execute(args []string) error {
	skinName := args[0]

	if o.account == "" {
		o.account = o.app.AccountManager().GetDefault()
	}

	token, err := o.app.AccountManager().GetAccountToken(o.account)
	if err != nil {
		return err
	}

	skin, err := o.app.SkinManager().GetSkin(skinName)
	if err != nil {
		return err
	}

	return skin.Apply(token)
}
