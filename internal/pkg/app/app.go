package app

import (
	"os"

	"github.com/mworzala/mc-cli/internal/pkg/account"
	"github.com/mworzala/mc-cli/internal/pkg/java"
	"github.com/mworzala/mc-cli/internal/pkg/platform"
	"github.com/spf13/cobra"
)

// App contains a wrapper for most of the actions performed by the CLI
type App struct {
	ConfigDir string

	cmd *cobra.Command

	accountManager account.Manager
	javaManager    java.Manager
}

func NewApp(cmd *cobra.Command) *App {
	configDir, err := platform.GetConfigDir()
	if err != nil {
		exitWithError(err)
	}

	return &App{
		ConfigDir: configDir,
		cmd:       cmd,
	}
}

func (a *App) AccountManager() account.Manager {
	if a.accountManager == nil {
		var err error
		a.accountManager, err = account.NewManager(a.ConfigDir)
		if err != nil {
			exitWithError(err)
		}
	}

	return a.accountManager
}

func (a *App) JavaManager() java.Manager {
	if a.javaManager == nil {
		var err error
		a.javaManager, err = java.NewManager(a.ConfigDir)
		if err != nil {
			exitWithError(err)
		}
	}

	return a.javaManager
}

func exitWithError(err error) {
	println("an error has occurred", err)
	os.Exit(1)
}
