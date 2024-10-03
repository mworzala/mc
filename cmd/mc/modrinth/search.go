package modrinth

import (
	"context"
	"os"
	"os/signal"
	"strings"

	"github.com/mworzala/mc/internal/pkg/cli"
	appModel "github.com/mworzala/mc/internal/pkg/cli/model"
	"github.com/mworzala/mc/internal/pkg/modrinth"
	"github.com/mworzala/mc/internal/pkg/modrinth/facet"
	"github.com/spf13/cobra"
)

type searchOpts struct {
	app *cli.App

	// Project type
	mod          bool
	modPack      bool
	resourcePack bool
	shader       bool

	// Sort
	//todo
}

func newSearchCmd(app *cli.App) *cobra.Command {
	var o searchOpts

	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search for projects on modrinth",
		Args: func(cmd *cobra.Command, args []string) error {
			o.app = app
			return o.validateArgs(cmd, args)
		},
		RunE: func(_ *cobra.Command, args []string) error {
			o.app = app
			return o.execute(args)
		},
	}

	cmd.Flags().BoolVar(&o.mod, "mod", false, "Show only mods")
	cmd.Flags().BoolVar(&o.modPack, "modpack", false, "Show only modpacks")
	cmd.Flags().BoolVar(&o.resourcePack, "resourcepack", false, "Show only resource packs")
	cmd.Flags().BoolVar(&o.resourcePack, "rp", false, "Show only resource packs")
	cmd.Flags().BoolVar(&o.shader, "shader", false, "Show only shaders")

	cmd.Flags().FlagUsages()

	return cmd
}

func (o *searchOpts) validateArgs(cmd *cobra.Command, args []string) (err error) {
	if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
		return err
	}

	// todo
	return nil
}

func (o *searchOpts) execute(args []string) error {
	// Validation function has done arg validation and option population

	client := modrinth.NewClient(o.app.Build.Version)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	query := strings.Join(args, " ")
	var facets facet.And
	if o.mod || o.modPack || o.resourcePack || o.shader {
		var filter facet.Or
		if o.mod {
			filter = append(filter, facet.Eq{facet.ProjectType, "mod"})
		}
		if o.modPack {
			filter = append(filter, facet.Eq{facet.ProjectType, "modpack"})
		}
		if o.resourcePack {
			filter = append(filter, facet.Eq{facet.ProjectType, "resourcepack"})
		}
		if o.shader {
			filter = append(filter, facet.Eq{facet.ProjectType, "shader"})
		}
		facets = append(facets, filter)
	}

	res, err := client.Search(ctx, modrinth.SearchRequest{
		Query:  query,
		Facets: facets,
	})
	if err != nil {
		return err
	}

	presentable := appModel.ModrinthSearchResult(*res)
	return o.app.Present(&presentable)
}
