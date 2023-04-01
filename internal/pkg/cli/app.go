package cli

import (
	"github.com/mworzala/mc-cli/internal/pkg/account"
	"github.com/mworzala/mc-cli/internal/pkg/cli/output"
	"github.com/mworzala/mc-cli/internal/pkg/game"
	"github.com/mworzala/mc-cli/internal/pkg/java"
	"github.com/mworzala/mc-cli/internal/pkg/platform"
	"github.com/mworzala/mc-cli/internal/pkg/profile"
)

type App struct {
	ConfigDir string
	Output    output.Format

	accountManager account.Manager
	javaManager    java.Manager
	versionManager *game.VersionManager
	profileManager profile.Manager
}

func NewApp() *App {
	configDir, err := platform.GetConfigDir()
	if err != nil {

	}

	return &App{
		ConfigDir: configDir,
		Output:    output.Format{Type: output.Default},
	}
}

func (a *App) AccountManager() account.Manager {
	if a.accountManager == nil {
		var err error
		a.accountManager, err = account.NewManager(a.ConfigDir)
		if err != nil {
			a.Fatal(err)
		}
	}

	return a.accountManager
}

func (a *App) JavaManager() java.Manager {
	if a.javaManager == nil {
		var err error
		a.javaManager, err = java.NewManager(a.ConfigDir)
		if err != nil {
			a.Fatal(err)
		}
	}

	return a.javaManager
}

func (a *App) VersionManager() *game.VersionManager {
	if a.versionManager == nil {
		var err error
		a.versionManager, err = game.NewVersionManager(a.ConfigDir)
		if err != nil {
			a.Fatal(err)
		}
	}

	return a.versionManager
}

func (a *App) ProfileManager() profile.Manager {
	if a.profileManager == nil {
		var err error
		a.profileManager, err = profile.NewManager(a.ConfigDir)
		if err != nil {
			a.Fatal(err)
		}
	}

	return a.profileManager
}
