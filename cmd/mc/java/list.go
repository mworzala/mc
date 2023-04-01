package java

import (
	appModel "github.com/mworzala/mc-cli/internal/pkg/app/model"
	"github.com/mworzala/mc-cli/internal/pkg/cli"
	"github.com/spf13/cobra"
)

type listJavaOpts struct {
	app *cli.App
}

func newListCmd(app *cli.App) *cobra.Command {
	var o listJavaOpts

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List discovered java installations",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			o.app = app
			return o.listInstallations()
		},
	}

	return cmd
}

func (o *listJavaOpts) listInstallations() error {
	javaManager := o.app.JavaManager()

	var result appModel.JavaInstallationList
	for _, name := range javaManager.Installations() {
		install := javaManager.GetInstallation(name)
		result = append(result, &appModel.JavaInstallation{
			Name:    install.Name,
			Path:    install.Path,
			Version: install.Version,
		})
	}

	return o.app.Present(result)
}
