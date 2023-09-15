package account

import (
	"fmt"

	"github.com/mworzala/mc/internal/pkg/account"
	"github.com/mworzala/mc/internal/pkg/cli"
	appModel "github.com/mworzala/mc/internal/pkg/cli/model"
	"github.com/mworzala/mc/internal/pkg/platform"
	"github.com/spf13/cobra"
)

type loginAccountOpts struct {
	app *cli.App

	accountType account.Type
	setDefault  bool
}

func newLoginCmd(app *cli.App) *cobra.Command {
	var o loginAccountOpts

	cmd := &cobra.Command{
		Use:        "login",
		Short:      "Sign into a new Minecraft account",
		ValidArgs:  []string{"microsoft", "mojang"},
		ArgAliases: []string{"mso", "minecraft", "mc"},
		Args:       cobra.MatchAll(cobra.MaximumNArgs(1), cobra.OnlyValidArgs),
		RunE: func(_ *cobra.Command, args []string) (err error) {
			o.app = app

			accountType := account.Microsoft
			if len(args) > 0 {
				accountType, err = account.ParseType(args[0])
				if err != nil {
					return
				}
			}
			o.accountType = accountType

			return o.handleLogin()
		},
	}

	cmd.Flags().BoolVar(&o.setDefault, "set-default", false, "Set the new account as the default")

	return cmd
}

func (o *loginAccountOpts) handleLogin() (err error) {
	accountManager := o.app.AccountManager()

	var acc *account.Account
	if o.accountType == account.Microsoft {
		acc, err = accountManager.LoginMicrosoft(func(verificationUrl, userCode string) {
			//todo should have a global flag to disable interactive elements like this
			_ = platform.OpenUrl(verificationUrl)
			_ = platform.WriteToClipboard(userCode)

			err := o.app.Present(&appModel.LoginPrompt{
				Url:  verificationUrl,
				Code: userCode,
			})
			if err != nil {
				o.app.Fatal(err)
			}
		})
	} else {
		return fmt.Errorf("mojang login currently not supported")
	}

	// Update the default in case the user requested to update the default,
	// or if there is no default set indicating that this is the first account
	if o.setDefault || accountManager.GetDefault() == "" {
		if err := accountManager.SetDefault(acc.UUID); err != nil {
			return fmt.Errorf("failed to set default account: %w", err)
		}
	}

	if err := accountManager.Save(); err != nil {
		return err
	}

	return o.app.Present(&appModel.Account{
		UUID:     acc.UUID,
		Username: acc.Profile.Username,
	})
}
