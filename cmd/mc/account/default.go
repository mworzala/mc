package account

import (
	"fmt"

	"github.com/mworzala/mc-cli/internal/pkg/account"
	"github.com/mworzala/mc-cli/internal/pkg/platform"
	"github.com/spf13/cobra"
)

//todo flags to emit the uuid instead of the name when getting the default
//todo if i add a json output for all the commands, output the whole profile

var defaultCmd = &cobra.Command{
	Use:     "default",
	Aliases: []string{"use"},
	Short:   "Get or set the default account",
	Long:    "abcd",
	Args:    cobra.MaximumNArgs(1),
	//todo override completions
	RunE: handleDefault,
}

func handleDefault(_ *cobra.Command, args []string) (err error) {
	dataDir, err := platform.GetConfigDir()
	if err != nil {
		return err
	}

	manager, err := account.NewManager(dataDir)
	if err != nil {
		return fmt.Errorf("failed to read accounts: %w", err)
	}

	if len(args) == 0 {
		return showDefaultAccount(manager)
	}
	return setDefaultAccount(manager, args[0])
}

func showDefaultAccount(m account.Manager) error {
	accountId := m.GetDefault()
	if accountId == "" {
		println("No accounts present, use `mc account login`") //todo better message (reminder to just search for println calls)
		return nil
	}

	acc := m.GetAccount(accountId, account.ModeUUID)
	if acc == nil {
		return fmt.Errorf("no account with default id present: %s", accountId)
	}

	println("Default account:", acc.Profile.Username)
	return nil
}

func setDefaultAccount(m account.Manager, newValue string) error {
	acc := m.GetAccount(newValue, account.ModeUUID|account.ModeName)
	if acc == nil {
		//todo error cases like this need to exit with code 1
		println("No account with name:", newValue)
	}

	if err := m.SetDefault(acc.UUID); err != nil {
		return err
	}

	if err := m.Save(); err != nil {
		return err
	}

	return nil
}
