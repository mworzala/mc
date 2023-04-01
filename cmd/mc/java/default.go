package java

import (
	"fmt"

	"github.com/mworzala/mc-cli/internal/pkg/app"

	"github.com/mworzala/mc-cli/internal/pkg/java"
	"github.com/spf13/cobra"
)

var defaultCmd = &cobra.Command{
	Use:     "default",
	Aliases: []string{"use"},
	Short:   "Get or set the default java installation",
	Args:    cobra.MaximumNArgs(1),
	RunE:    handleDefault,
}

func handleDefault(cmd *cobra.Command, args []string) (err error) {
	a := app.NewApp(cmd)

	if len(args) == 0 {
		return showDefaultInstallation(a.JavaManager())
	}
	return setDefaultInstallation(a.JavaManager(), args[0])
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
