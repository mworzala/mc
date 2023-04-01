package account

import (
	"fmt"

	"github.com/mworzala/mc-cli/internal/pkg/account"
	"github.com/mworzala/mc-cli/internal/pkg/platform"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:        "login",
	Short:      "todo",
	Long:       "abcd",
	ValidArgs:  []string{"microsoft", "mojang"},
	ArgAliases: []string{"mso", "minecraft", "mc"},
	Args:       cobra.MatchAll(cobra.MaximumNArgs(1), cobra.OnlyValidArgs),
	RunE:       handleLogin,
}

var (
	setDefaultAccountFlag bool
)

func init() {
	loginCmd.Flags().BoolVar(&setDefaultAccountFlag, "set-default", false, "Set the new account as the default")
}

func handleLogin(_ *cobra.Command, args []string) (err error) {
	accountType := account.Microsoft
	if len(args) > 0 {
		var ok bool
		if accountType, ok = account.TypeFromString(args[0]); !ok {
			return fmt.Errorf("unknown account type: %s", args[0])
		}
	}

	dataDir, err := platform.GetConfigDir()
	if err != nil {
		return err
	}

	manager, err := account.NewManager(dataDir)
	if err != nil {
		return fmt.Errorf("failed to read accounts: %w", err)
	}

	var acc *account.Account
	if accountType == account.Microsoft {
		acc, err = manager.LoginMicrosoft(func(verificationUrl, userCode string) {
			_ = platform.OpenUrl(verificationUrl)
			_ = platform.WriteToClipboard(userCode)
			//todo better message
			println(verificationUrl, userCode)
		})
	} else {
		return fmt.Errorf("mojang login currently not supported")
	}

	// Update the default in case the user requested to update the default,
	// or if there is no default set indicating that this is the first account
	if setDefaultAccountFlag || manager.GetDefault() == "" {
		if err := manager.SetDefault(acc.UUID); err != nil {
			return fmt.Errorf("failed to set default account: %w", err)
		}
	}

	if err := manager.Save(); err != nil {
		return err
	}

	fmt.Printf("🎉 Signed in as %s\n", acc.Profile.Username)
	return
}
