package skin

import (
	"errors"
	"fmt"
	"slices"

	"github.com/google/uuid"
	"github.com/mworzala/mc/internal/pkg/cli"
	"github.com/spf13/cobra"
)

type addSkinOpts struct {
	app *cli.App

	account string
	variant string
	cape    string
	name    string
	apply   bool
}

var (
	ErrInvalidType    = errors.New("invalid type")
	ErrInvalidVariant = errors.New("invalid variant")

	validVariants = []string{"classic", "slim"}
)

func newAddCmd(app *cli.App) *cobra.Command {
	var o addSkinOpts

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a skin to your list",
		Args: func(cmd *cobra.Command, args []string) error {
			o.app = app
			return o.validateArgs(cmd, args)
		},
		RunE: func(_ *cobra.Command, args []string) error {
			o.app = app
			return o.execute(args)
		},
	}

	cmd.Flags().StringVar(&o.account, "account", "", "Account to use")
	cmd.Flags().StringVar(&o.variant, "variant", "classic", "Skin variant [classic/slim]")
	cmd.Flags().StringVar(&o.cape, "cape", "", "Cape name, 'none' to remove")
	cmd.Flags().BoolVar(&o.apply, "apply", false, "Apply the skin")
	cmd.Flags().BoolVar(&o.apply, "set", false, "Apply the skin")

	cmd.Flags().FlagUsages()

	return cmd
}

func (o *addSkinOpts) validateArgs(cmd *cobra.Command, args []string) (err error) {
	if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
		return err
	}

	if !slices.Contains(validVariants, o.variant) {
		return ErrInvalidVariant
	}

	return nil
}

func (o *addSkinOpts) execute(args []string) error {
	if len(args) > 1 {
		o.name = args[1]
	}

	if o.name == "" {
		o.name = uuid.New().String()
	}

	if o.account == "" {
		o.account = o.app.AccountManager().GetDefault()
	}

	token, err := o.app.AccountManager().GetAccountToken(o.account)
	if err != nil {
		return err
	}

	info, err := o.app.SkinManager().GetProfileInformation(token)
	if err != nil {
		return err
	}

	if o.cape == "" {
		for _, cape := range info.Capes {
			if cape.State == "ACTIVE" {
				o.cape = cape.ID
			}
		}
	} else if o.cape != "none" {
		for _, cape := range info.Capes {
			if cape.Alias == o.cape {
				o.cape = cape.ID
			}
		}
	}

	skinData := args[0]

	skin, err := o.app.SkinManager().CreateSkin(o.name, o.variant, skinData, o.cape)
	if err != nil {
		return err
	}

	if o.apply {

		err = skin.Apply(token)
		if err != nil {
			return err
		}
	}

	fmt.Printf("skin %s with cape %s was added to the list", skin.Name, skin.Cape)

	return o.app.SkinManager().Save()

}
