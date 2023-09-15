package account

import (
	"errors"
	"fmt"

	"github.com/mworzala/mc/internal/pkg/cli"
	appModel "github.com/mworzala/mc/internal/pkg/cli/model"

	"github.com/spf13/cobra"
)

type defaultAccountOpts struct {
	app *cli.App
}

func newDefaultCmd(app *cli.App) *cobra.Command {
	var o defaultAccountOpts

	cmd := &cobra.Command{
		Use:     "default",
		Aliases: []string{"use"},
		Short:   "Manage the default account",
		Args:    cobra.MaximumNArgs(1),
		//todo override completions
		RunE: func(_ *cobra.Command, args []string) error {
			o.app = app

			if len(args) == 0 {
				return o.getDefault()
			}
			return o.setDefault(args)
		},
	}

	return cmd
}

func (o *defaultAccountOpts) getDefault() error {
	accountManager := o.app.AccountManager()
	accounts := accountManager.Accounts()
	if len(accounts) == 0 {
		return errors.New("no accounts configured. try 'mc account login'")
	}

	// Read default account
	acc := accountManager.GetAccount(accountManager.GetDefault())
	if acc == nil {
		// In this case the default has been misconfigured because we know there is at least one account, but yet the default does not exist.
		// Correct the issue by resetting the default account to the first known account
		if err := accountManager.SetDefault(accounts[0]); err != nil {
			return err
		}
		if err := accountManager.Save(); err != nil {
			return err
		}

		// Now we know this is a safe call
		acc = accountManager.GetAccount(accountManager.GetDefault())
	}

	return o.app.Present(&appModel.Account{
		UUID:     acc.UUID,
		Username: acc.Profile.Username,
	})
}

func (o *defaultAccountOpts) setDefault(args []string) error {
	accountManager := o.app.AccountManager()

	// Validate new account
	acc := accountManager.GetAccount(args[0])
	if acc == nil {
		return fmt.Errorf("no such account: %s", args[0])
	}

	// Update and save
	if err := accountManager.SetDefault(acc.UUID); err != nil {
		return err
	}
	return accountManager.Save()
}
