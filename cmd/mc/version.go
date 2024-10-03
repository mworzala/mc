package mc

import (
	"fmt"
	"runtime"

	"github.com/mworzala/mc/internal/pkg/cli"
	appModel "github.com/mworzala/mc/internal/pkg/cli/model"
	"github.com/spf13/cobra"
)

func newVersionCmd(app *cli.App) *cobra.Command {

	cmd := &cobra.Command{
		Use:    "version",
		Hidden: true,
		Args:   cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			state := "clean"
			if app.Build.Modified {
				state = "modified"
			}

			latestRelease, latestSnapshot, err := app.VersionManager().LatestVersions()
			if err != nil {
				return err
			}

			return app.Present(&appModel.Version{
				Tag:      app.Build.Version,
				Commit:   app.Build.Commit,
				State:    state,
				Date:     app.Build.Date,
				Go:       runtime.Version(),
				Platform: fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
				Game: appModel.GameVersion{
					Release:  latestRelease,
					Snapshot: latestSnapshot,
				},
			})
		},
	}

	return cmd
}
