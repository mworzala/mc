package mc

import (
	"fmt"
	"runtime"

	"github.com/mworzala/mc-cli/internal/pkg/cli"
	appModel "github.com/mworzala/mc-cli/internal/pkg/cli/model"
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

			return app.Present(&appModel.Version{
				Tag:      app.Build.Version,
				Commit:   app.Build.Commit,
				State:    state,
				Date:     app.Build.Date,
				Go:       runtime.Version(),
				Platform: fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
			})
		},
	}

	return cmd
}
