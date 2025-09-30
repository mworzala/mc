package mod

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/mworzala/mc/internal/pkg/cli"
	"github.com/mworzala/mc/internal/pkg/modrinth"
	"github.com/mworzala/mc/internal/pkg/profile"
	"github.com/spf13/cobra"
)

var (
	errNotModded  = errors.New("profile is not a modded profile")
	errNoVersions = errors.New("no versions found")
)

type modInstallOpts struct {
	app       *cli.App
	versionId string
}

func newModInstallCmd(app *cli.App) *cobra.Command {
	var o modInstallOpts

	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install a mod",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errNoProfile
			}
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			versionProvided := cmd.Flags().Changed("version-id")

			// a mod id isn't required if a version id is specified
			if !versionProvided && len(args) < 2 {
				return errNoModId
			}
			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			o.app = app
			return o.installMod(args)
		},
	}

	cmd.Flags().StringVar(&o.versionId, "version-id", "", "Install from a specific version id")

	return cmd
}

func (o *modInstallOpts) installMod(args []string) error {
	profileManager := o.app.ProfileManager()
	profileName := args[0]
	p, err := profileManager.GetProfile(profileName)
	if err != nil {
		return fmt.Errorf("%w: %s", err, profileName)
	}

	if p.Type != profile.Fabric {
		return fmt.Errorf("%w: %s", errNotModded, profileName)
	}

	client := modrinth.NewClient(o.app.Build.Version)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	gameVersions := []string{p.GameVersion()}
	// todo: when supporting forge, actually pass in the type
	loaders := []string{"fabric"}

	var version modrinth.ProjectVersion

	if o.versionId != "" {
		v, err := client.GetProjectVersion(ctx, o.versionId)
		if err != nil {
			return fmt.Errorf("%w: %s", err, o.versionId)
		}
		version = *v
	} else {
		modId := args[1]

		// todo: direct version id support
		versions, err := client.GetProjectVersions(ctx, modId, gameVersions, loaders)
		if err != nil {
			return fmt.Errorf("%w: %s", err, modId)
		}

		if len(*versions) < 1 {
			return fmt.Errorf("%w: %s", errNoVersions, modId)
		}

		version = (*versions)[0]
	}

	deps, err := gatherDependencies(client, ctx, gameVersions, loaders, version)
	if err != nil {
		return err
	}

	var mods []profile.UnsavedMod
	mods = append(mods, profile.UnsavedMod{
		Id:        version.ProjectId,
		VersionId: version.Id,
		File:      *findPrimaryFile(version),
		Provider:  profile.Modrinth, // curseforge isn't supported yet but in the future it may be
	})
	for _, dep := range deps {
		mods = append(mods, profile.UnsavedMod{
			Id:        dep.ProjectId,
			VersionId: dep.Id,
			File:      *findPrimaryFile(dep),
			Provider:  profile.Modrinth, // curseforge isn't supported yet but in the future it may be
		})
	}

	for _, mod := range mods {
		if !o.app.Config.NonInteractive {
			fmt.Printf("Downloading and Installing %s\n", mod.File.Filename)
		}
		err = p.InstallMod(client, ctx, mod)
		if err != nil {
			if errors.Is(err, profile.ErrModAlreadyInstalled) {
				if !o.app.Config.NonInteractive {
					fmt.Printf("Skipping %s as its already installed\n", mod.File.Filename)
				}
				continue
			}
			return fmt.Errorf("%w: %s", err, mod.File.Filename)
		}
		if !o.app.Config.NonInteractive {
			fmt.Printf("Done Installing %s\n", mod.File.Filename)
		}
	}

	if err = p.SaveMods(); err != nil {
		return err
	}

	return nil
}

func gatherDependencies(client *modrinth.Client, ctx context.Context, gameVersions []string, loaders []string, version modrinth.ProjectVersion) ([]modrinth.ProjectVersion, error) {
	var depVersions []modrinth.ProjectVersion
	for _, dependency := range version.Dependencies {
		if dependency.DependencyType != modrinth.DependencyRequired {
			continue
		}

		if dependency.VersionId != nil {
			v, err := client.GetProjectVersion(ctx, *dependency.VersionId)
			if err != nil {
				return nil, err
			}
			depVersions = append(depVersions, *v)

			deeperDeps, err := gatherDependencies(client, ctx, gameVersions, loaders, *v)
			if err != nil {
				return nil, err
			}
			depVersions = append(depVersions, deeperDeps...)

		} else if dependency.ProjectId != nil {
			versions, err := client.GetProjectVersions(ctx, *dependency.ProjectId, gameVersions, loaders)
			if err != nil {
				return nil, err
			}

			if len(*versions) < 1 {
				return nil, fmt.Errorf("%w: %s", errNoVersions, *dependency.ProjectId)
			}

			v := (*versions)[0]

			depVersions = append(depVersions, v)

			deeperDeps, err := gatherDependencies(client, ctx, gameVersions, loaders, v)
			if err != nil {
				return nil, err
			}
			depVersions = append(depVersions, deeperDeps...)

		}

	}

	return depVersions, nil
}

func findPrimaryFile(version modrinth.ProjectVersion) *modrinth.VersionFile {
	for _, file := range version.Files {
		if file.Primary {
			return &file
		}
	}

	// Modrinth's docs say "If there are not any primary files, it can be inferred that the first file is the primary one."
	return &version.Files[0]
}
