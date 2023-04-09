package mc

import (
	"fmt"

	"github.com/mworzala/mc-cli/internal/pkg/cli"

	"github.com/mworzala/mc-cli/internal/pkg/account"
	"github.com/mworzala/mc-cli/internal/pkg/game/launch"
	"github.com/spf13/cobra"
)

type launchOpts struct {
	app *cli.App
}

func newLaunchCmd(app *cli.App) *cobra.Command {
	var o launchOpts

	cmd := &cobra.Command{
		Use:     "launch",
		Aliases: []string{"run"},
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			o.app = app
			return o.launch(args)
		},
	}

	return cmd
}

func (o *launchOpts) launch(args []string) error {

	profileManager := o.app.ProfileManager()
	p, err := profileManager.GetProfile(args[0])
	if err != nil {
		return fmt.Errorf("%w: %s", err, args[0])
	}

	accountManager := o.app.AccountManager()
	acc := accountManager.GetAccount(accountManager.GetDefault(), account.ModeUUID)
	if acc == nil {
		return fmt.Errorf("no default account is set")
	}

	javaManager := o.app.JavaManager()
	javaInstall := javaManager.GetInstallation(javaManager.GetDefault())
	if javaInstall == nil {
		return fmt.Errorf("no default java installation is set")
	}

	return launch.LaunchProfile(o.app.ConfigDir, p, acc, javaInstall)

}
