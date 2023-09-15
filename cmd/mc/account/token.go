package account

import (
	"fmt"

	"github.com/mworzala/mc-cli/internal/pkg/cli"
	appModel "github.com/mworzala/mc-cli/internal/pkg/cli/model"
	"github.com/spf13/cobra"
)

type tokenAccountOpts struct {
	app *cli.App
}

func newTokenCmd(app *cli.App) *cobra.Command {
	var o tokenAccountOpts

	cmd := &cobra.Command{
		Use:   "token",
		Short: "Get a Minecraft access token",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) (err error) {
			o.app = app

			return o.handleToken(args)
		},
	}

	return cmd
}

func (o *tokenAccountOpts) handleToken(args []string) (err error) {
	accountManager := o.app.AccountManager()

	var account string
	if len(args) > 0 {
		account = args[0]
	} else {
		account = accountManager.GetDefault()
	}

	acc := accountManager.GetAccount(account)
	if acc == nil {
		return fmt.Errorf("no such account")
	}

	token, err := accountManager.GetAccountToken(acc.UUID)
	if err != nil {
		return err //todo better error messages
	}

	return o.app.Present(&appModel.AccessToken{
		Username: acc.Profile.Username,
		UUID:     acc.UUID,
		Token:    token,
	})
}
