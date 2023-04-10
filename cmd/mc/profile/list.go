package profile

import (
	"github.com/mworzala/mc-cli/internal/pkg/cli"
	appModel "github.com/mworzala/mc-cli/internal/pkg/cli/model"
	"github.com/spf13/cobra"
)

type listProfilesOpts struct {
	app *cli.App
}

func newListCmd(app *cli.App) *cobra.Command {
	var o listProfilesOpts

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List installed profiles",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			o.app = app
			return o.listProfiles()
		},
	}

	return cmd
}

func (o *listProfilesOpts) listProfiles() error {
	profileManager := o.app.ProfileManager()

	var result appModel.ProfileList
	for _, name := range profileManager.Profiles() {
		profile, _ := profileManager.GetProfile(name) // Ignore error since we just got the list of names
		result = append(result, &appModel.Profile{
			Name:      profile.Name,
			Directory: profile.Directory,
			Type:      appModel.ProfileTypes[profile.Type],
			Version:   profile.Version,
		})
	}

	return o.app.Present(result)
}
