package skin

import (
	"github.com/mworzala/mc/internal/pkg/cli"
	appModel "github.com/mworzala/mc/internal/pkg/cli/model"
	"github.com/spf13/cobra"
)

type listSkinsOpts struct {
	app *cli.App
}

func newListCmd(app *cli.App) *cobra.Command {
	var o listSkinsOpts

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List saved skins",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, args []string) error {
			o.app = app
			return o.listSkins()
		},
	}

	return cmd
}

func (o *listSkinsOpts) listSkins() error {
	skinManager := o.app.SkinManager()

	var result appModel.SkinList
	for _, skin := range skinManager.Skins() {
		result = append(result, &appModel.Skin{
			Name:     skin.Name,
			Modified: skin.AddedDate,
		})
	}

	return o.app.Present(result)
}
