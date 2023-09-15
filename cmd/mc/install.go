package mc

import (
	"errors"
	"fmt"

	"github.com/mworzala/mc-cli/internal/pkg/cli"
	"github.com/mworzala/mc-cli/internal/pkg/game"
	"github.com/mworzala/mc-cli/internal/pkg/game/install"
	gameModel "github.com/mworzala/mc-cli/internal/pkg/game/model"
	"github.com/mworzala/mc-cli/internal/pkg/profile"
	"github.com/spf13/cobra"
)

//todo tab completions

type installOpts struct {
	app *cli.App

	version *gameModel.VersionInfo

	fabric       bool
	fabricLoader string
}

func newInstallCmd(app *cli.App) *cobra.Command {
	var o installOpts

	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install a new Minecraft version",
		Args: func(cmd *cobra.Command, args []string) error {
			o.app = app
			return o.validateArgs(cmd, args)
		},
		RunE: func(_ *cobra.Command, args []string) error {
			o.app = app
			return o.installSelected(args)
		},
	}

	cmd.Flags().BoolVar(&o.fabric, "fabric", false, "Install fabric mod loader")
	cmd.Flags().StringVar(&o.fabricLoader, "loader", "", "Fabric loader version, ignored without --fabric")

	return cmd
}

func (o *installOpts) validateArgs(cmd *cobra.Command, args []string) (err error) {
	if err := cobra.RangeArgs(1, 2)(cmd, args); err != nil {
		return err
	}

	// Validate version arg (0)
	versionManager := o.app.VersionManager()
	if o.version, err = versionManager.FindVanilla(args[0]); errors.Is(err, game.ErrUnknownVersion) {
		return fmt.Errorf("%w: %s", err, args[0])
	}

	// Validate name arg (1, optional)
	if len(args) > 1 && !profile.IsValidName(args[1]) {
		return profile.ErrInvalidName
	}

	// Validate flag fabric (and loader)
	// If fabric is present ensure the selected version is supported, if loader is specified use that
	if o.fabric {
		if o.fabricLoader == "" {
			o.fabricLoader = versionManager.DefaultFabricLoader()
		}

		o.version, err = versionManager.FindFabric(args[0], o.fabricLoader)
		if errors.Is(err, game.ErrUnknownFabricVersion) {
			return fmt.Errorf("%w: %s", err, args[0])
		}
		if errors.Is(err, game.ErrUnknownFabricLoader) {
			return fmt.Errorf("%w: %s", err, o.fabricLoader)
		}
	}

	return nil
}

func (o *installOpts) installSelected(args []string) error {
	// Validation function has done arg validation and option population

	// Install the selected version
	versionManager := o.app.VersionManager()
	installer := install.NewInstaller(o.app.ConfigDir, versionManager.FindVanilla)
	if err := installer.Install(o.version); err != nil {
		return fmt.Errorf("installation failed: %w", err)
	}

	// Initialize profile
	profileManager := o.app.ProfileManager()

	profileName := args[0]
	if o.fabric {
		profileName = fmt.Sprintf("%s-fabric", profileName)
	}
	if len(args) > 1 {
		profileName = args[1]
	}
	p, err := profileManager.CreateProfile(profileName)
	if err != nil {
		return err
	}

	p.Type = profile.Vanilla
	if o.fabric {
		p.Type = profile.Fabric
	}
	p.Version = o.version.Id

	if err := profileManager.Save(); err != nil {
		return err
	}

	println("installed", o.version.Id)
	return nil
}
