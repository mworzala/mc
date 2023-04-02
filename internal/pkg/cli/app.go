package cli

import (
	"github.com/mworzala/mc-cli/internal/pkg/account"
	"github.com/mworzala/mc-cli/internal/pkg/cli/output"
	"github.com/mworzala/mc-cli/internal/pkg/game"
	"github.com/mworzala/mc-cli/internal/pkg/java"
	"github.com/mworzala/mc-cli/internal/pkg/platform"
	"github.com/mworzala/mc-cli/internal/pkg/profile"
)

type BuildInfo struct {
	Version  string
	Commit   string
	Date     string
	Modified bool
}

type App struct {
	Build BuildInfo

	ConfigDir string
	Output    output.Format

	accountManager account.Manager
	javaManager    java.Manager
	versionManager *game.VersionManager
	profileManager profile.Manager
}

func NewApp(build BuildInfo) *App {
	a := &App{Build: build}

	var err error
	if a.ConfigDir, err = platform.GetConfigDir(build.Version == "dev"); err != nil {
		a.Fatal(err)
	}

	a.Output = output.Format{Type: output.Default}
	return a
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
