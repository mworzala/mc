package java

import (
	"fmt"

	"github.com/mworzala/mc-cli/internal/pkg/java"
	"github.com/mworzala/mc-cli/internal/pkg/platform"
	"github.com/spf13/cobra"
)

var defaultCmd = &cobra.Command{
	Use:     "default",
	Aliases: []string{"use"},
	Short:   "Get or set the default java installation",
	Args:    cobra.MaximumNArgs(1),
	RunE:    handleDefault,
}

func handleDefault(_ *cobra.Command, args []string) (err error) {
	dataDir, err := platform.GetConfigDir()
	if err != nil {
		return err
	}

	manager, err := java.NewManager(dataDir)
	if err != nil {
		return fmt.Errorf("failed to read accounts: %w", err)
	}

	if len(args) == 0 {
		return showDefaultInstallation(manager)
	}
	return setDefaultInstallation(manager, args[0])
}

func showDefaultInstallation(m java.Manager) error {
	javaInstall := m.GetDefault()
	if javaInstall == "" {
		println("No java installations")
		return nil
	}

	install := m.GetInstallation(javaInstall)
	if install == nil {
		return fmt.Errorf("no java installation with default name present: %s", javaInstall)
	}

	println("Default java:", install.Path)
	return nil
}

func setDefaultInstallation(m java.Manager, newValue string) error {
	install := m.GetInstallation(newValue)
	if install == nil {
		//todo error cases like this need to exit with code 1
		println("No installation with name:", newValue)
	}

	if err := m.SetDefault(install.Name); err != nil {
		return err
	}

	if err := m.Save(); err != nil {
		return err
	}

	return nil
}
