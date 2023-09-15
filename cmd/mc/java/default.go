package java

import (
	"errors"
	"fmt"

	"github.com/mworzala/mc/internal/pkg/cli"
	appModel "github.com/mworzala/mc/internal/pkg/cli/model"
	"github.com/spf13/cobra"
)

type defaultJavaOpts struct {
	app *cli.App
}

func newDefaultCmd(app *cli.App) *cobra.Command {
	var o defaultJavaOpts

	cmd := &cobra.Command{
		Use:     "default",
		Aliases: []string{"use"},
		Short:   "Manage the default java installation",
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			o.app = app

			if len(args) == 0 {
				return o.getDefault()
			}
			return o.setDefault(args)
		},
	}

	return cmd
}

func (o *defaultJavaOpts) getDefault() error {
	javaManager := o.app.JavaManager()
	installs := javaManager.Installations()
	if len(installs) == 0 {
		return errors.New("no java installations found")
	}

	// Read default installation
	install := javaManager.GetInstallation(javaManager.GetDefault())
	if install == nil {
		// In this case the default has been misconfigured because we know there is at least one installation, but yet the default does not exist.
		// Correct the issue by resetting the default installation to the first known installation
		if err := javaManager.SetDefault(installs[0]); err != nil {
			return err
		}
		if err := javaManager.Save(); err != nil {
			return err
		}

		// Now we know this is a safe call
		install = javaManager.GetInstallation(javaManager.GetDefault())
	}

	return o.app.Present(&appModel.JavaInstallation{
		Name:    install.Name,
		Path:    install.Path,
		Version: install.Version,
	})
}

func (o *defaultJavaOpts) setDefault(args []string) error {
	javaManager := o.app.JavaManager()

	// Validate new installation existence
	install := javaManager.GetInstallation(args[0])
	if install == nil {
		return fmt.Errorf("no such java installation: %s", args[0])
	}

	// Update and save
	if err := javaManager.SetDefault(install.Name); err != nil {
		return err
	}
	return javaManager.Save()
}
