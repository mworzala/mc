package main

import (
	"fmt"

	"github.com/mworzala/mc-cli/internal/pkg/account"
	"github.com/mworzala/mc-cli/internal/pkg/game/launch"
	"github.com/mworzala/mc-cli/internal/pkg/java"
	"github.com/mworzala/mc-cli/internal/pkg/platform"
	"github.com/mworzala/mc-cli/internal/pkg/profile"
	"github.com/spf13/cobra"
)

var launchCmd = &cobra.Command{
	Aliases: []string{"launch", "run"},
	RunE:    handleLaunch,
}

func handleLaunch(_ *cobra.Command, args []string) error {

	dataDir, err := platform.GetConfigDir()
	if err != nil {
		return err
	}

	profileManager, err := profile.NewManager(dataDir)
	if err != nil {
		return err
	}

	p := profileManager.GetProfile("1.19.2")
	if p == nil {
		panic("not installed")
	}

	accountManager, err := account.NewManager(dataDir)
	if err != nil {
		return err
	}

	javaManager, err := java.NewManager(dataDir)
	if err != nil {
		return err
	}
	defaultJava := javaManager.GetDefault()
	if defaultJava == "" {
		return fmt.Errorf("a default java must be configured.")
	}
	javaInstall := javaManager.GetInstallation(defaultJava)

	acc := accountManager.GetAccount(accountManager.GetDefault(), account.ModeUUID)

	return launch.LaunchProfile(dataDir, p, acc, javaInstall)

}
