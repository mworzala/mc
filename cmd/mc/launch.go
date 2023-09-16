package mc

import (
	"fmt"
	"os"

	"github.com/mworzala/mc/internal/pkg/java"

	"github.com/mworzala/mc/internal/pkg/cli"

	"github.com/mworzala/mc/internal/pkg/game/launch"
	"github.com/spf13/cobra"
)

type launchOpts struct {
	app *cli.App

	// Quick play actions
	quickPlaySingleplayer string
	quickPlayMultiplayer  string
	quickPlayRealms       string

	tail bool
}

func newLaunchCmd(app *cli.App) *cobra.Command {
	var o launchOpts

	cmd := &cobra.Command{
		Use:     "launch",
		Short:   "Launch a profile (Minecraft installation)",
		Aliases: []string{"run"},
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			o.app = app
			return o.launch(args)
		},
	}

	cmd.Flags().StringVarP(&o.quickPlaySingleplayer, "world", "", "", "launch into the server (1.20+)")
	cmd.Flags().StringVarP(&o.quickPlayMultiplayer, "server", "", "", "launch into a world (1.20+)")
	cmd.Flags().StringVarP(&o.quickPlayRealms, "realm", "", "", "launch into a realm (1.20+)")
	cmd.MarkFlagsMutuallyExclusive("world", "server", "realm")

	cmd.Flags().BoolVarP(&o.tail, "tail", "t", false, "attach the game stdout to the process")

	return cmd
}

func (o *launchOpts) launch(args []string) error {

	profileManager := o.app.ProfileManager()
	p, err := profileManager.GetProfile(args[0])
	if err != nil {
		return fmt.Errorf("%w: %s", err, args[0])
	}

	accountManager := o.app.AccountManager()
	acc := accountManager.GetAccount(accountManager.GetDefault())
	if acc == nil {
		return fmt.Errorf("no default account is set")
	}
	accessToken, err := accountManager.GetAccountToken(acc.UUID)
	if err != nil {
		return err
	}

	javaManager := o.app.JavaManager()
	var javaInstall *java.Installation
	if p.Config().Java != "" {
		javaInstall = javaManager.GetInstallation(p.Config().Java)
		if javaInstall == nil {
			_, _ = fmt.Fprintf(os.Stderr, "warning: configured java installation '%s' not found, using default\n", p.Config().Java)
		}
	}
	// If still unset (or was invalid), use the default
	if javaInstall == nil {
		javaInstall = javaManager.GetInstallation(javaManager.GetDefault())
		if javaInstall == nil {
			return fmt.Errorf("no default java installation is set")
		}
	}

	var quickPlay *launch.QuickPlay
	if o.quickPlaySingleplayer != "" {
		quickPlay = &launch.QuickPlay{
			Type: launch.QuickPlaySingleplayer,
			Id:   o.quickPlaySingleplayer,
		}
	} else if o.quickPlayMultiplayer != "" {
		quickPlay = &launch.QuickPlay{
			Type: launch.QuickPlayMultiplayer,
			Id:   o.quickPlayMultiplayer,
		}
	} else if o.quickPlayRealms != "" {
		quickPlay = &launch.QuickPlay{
			Type: launch.QuickPlayRealms,
			Id:   o.quickPlayRealms,
		}
	}

	return launch.LaunchProfile(o.app.ConfigDir, p, acc, accessToken, javaInstall, o.tail, quickPlay)
}
