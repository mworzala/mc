package cli

import (
	"errors"
	"strings"

	"github.com/mworzala/mc/internal/pkg/account"
	"github.com/mworzala/mc/internal/pkg/cli/output"
	"github.com/mworzala/mc/internal/pkg/config"
	"github.com/mworzala/mc/internal/pkg/game"
	"github.com/mworzala/mc/internal/pkg/java"
	"github.com/mworzala/mc/internal/pkg/platform"
	"github.com/mworzala/mc/internal/pkg/profile"
	"github.com/spf13/viper"
)

type BuildInfo struct {
	Version  string
	Commit   string
	Date     string
	Modified bool
	Source   bool
}

type App struct {
	Build BuildInfo

	ConfigDir string
	Output    output.Format //todo migrate
	Config    *config.Config

	accountManager account.Manager
	javaManager    java.Manager
	versionManager *game.VersionManager
	profileManager profile.Manager
	gameManager    game.Manager
}

func NewApp(build BuildInfo) *App {
	a := &App{Build: build}

	var err error
	if a.ConfigDir, err = platform.GetConfigDir(build.Source); err != nil {
		a.Fatal(err)
	}

	a.readConfig()

	a.Output = output.Format{Type: output.Default}
	return a
}

func (a *App) readConfig() {
	v := viper.New()
	v.SetConfigType("toml")
	v.SetConfigName("config")
	v.AddConfigPath(a.ConfigDir)

	// Allow config options to be overridden by env vars
	v.SetEnvPrefix("MC")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Set default values
	v.SetDefault("use_system_keyring", true)

	err := v.ReadInConfig()
	if err != nil && !errors.As(err, &viper.ConfigFileNotFoundError{}) {
		// For some reason viper doesnt implement Is() for ConfigFileNotFoundError
		a.Fatal(err)
	}
	a.Config = &config.Config{}
	if err := v.Unmarshal(a.Config); err != nil {
		a.Fatal(err)
	}
}

func (a *App) AccountManager() account.Manager {
	if a.accountManager == nil {
		var err error
		a.accountManager, err = account.NewManager(a.ConfigDir, a.Config)
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

func (a *App) GameManager() game.Manager {
	if a.gameManager == nil {
		var err error
		a.gameManager, err = game.NewManager(a.ConfigDir)
		if err != nil {
			a.Fatal(err)
		}
	}

	return a.gameManager
}
