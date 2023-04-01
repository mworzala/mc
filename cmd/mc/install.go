package main

import (
	"fmt"

	"github.com/mworzala/mc-cli/internal/pkg/game"
	"github.com/mworzala/mc-cli/internal/pkg/platform"
	"github.com/mworzala/mc-cli/internal/pkg/profile"
	"github.com/spf13/cobra"
)

//todo tab completions

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install a new Minecraft version and/or mod loader)",
	Args:  validateInstallArgs,
	RunE:  handleInstall,
}

var (
	// Flags

	// Filled in during argument validation
	versionManager *game.VersionManager
	targetVersion  *game.VersionInfo
)

func validateInstallArgs(cmd *cobra.Command, args []string) error {
	println("validate args")
	if err := cobra.RangeArgs(1, 2)(cmd, args); err != nil {
		return err
	}

	// Validate version arg (0)
	dataDir, err := platform.GetConfigDir()
	if err != nil {
		return err
	}
	versionManager, err = game.NewVersionManager(dataDir)
	if err != nil {
		return err
	}
	targetVersion = versionManager.FindVersionByName(args[0])
	if targetVersion == nil {
		return fmt.Errorf("no known version named %s", args[0])
	}

	// Validate name arg (1, optional)
	if len(args) > 1 {
		profileName := args[1]
		if !profile.IsValidName(profileName) {
			return fmt.Errorf("invalid profile name")
		}
	}

	return nil
}

func handleInstall(_ *cobra.Command, args []string) error {
	// The validation function has already initialized versionManager & targetVersion

	dataDir, err := platform.GetConfigDir()
	if err != nil {
		return err
	}

	// Install version
	if err := game.InstallVersion(dataDir, targetVersion); err != nil {
		return err
	}

	// Initialize profile
	profileManager, err := profile.NewManager(dataDir)
	if err != nil {
		return err
	}

	profileName := args[0]
	if len(args) > 1 {
		profileName = args[1]
	}
	p, err := profileManager.CreateProfile(profileName)
	if err != nil {
		return err
	}

	p.Type = profile.Vanilla
	p.Version = targetVersion.Id
	println("installed minecraft", targetVersion.Id)

	if err := profileManager.Save(); err != nil {
		return err
	}

	return nil
}
