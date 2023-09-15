package java

import (
	"fmt"

	"github.com/mworzala/mc-cli/internal/pkg/cli"
	"github.com/mworzala/mc-cli/internal/pkg/cli/model"
	"github.com/spf13/cobra"
)

type discoverJavaOpts struct {
	app *cli.App

	setDefault bool
}

func newDiscoverCommand(app *cli.App) *cobra.Command {
	var o discoverJavaOpts

	cmd := &cobra.Command{
		Use:   "discover",
		Short: "Discover a new Java installation",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			o.app = app
			return o.discover(args)
		},
	}

	cmd.Flags().BoolVar(&o.setDefault, "set-default", false, "Set the new installation as the default")

	return cmd
}

func (o *discoverJavaOpts) discover(args []string) error {
	javaManager := o.app.JavaManager()
	install, err := javaManager.Discover(args[0])
	if err != nil {
		return fmt.Errorf("java discovery failed: %w", err)
	}

	// Update default if there is not one, or the flag was set
	if o.setDefault || javaManager.GetDefault() == "" {
		if err := javaManager.SetDefault(install.Name); err != nil {
			return err
		}
	}

	if err := javaManager.Save(); err != nil {
		return err
	}

	return o.app.Present(&model.JavaInstallation{
		Name:    install.Name,
		Path:    install.Path,
		Version: install.Version,
	})
}
