package app

import (
	"encoding/json"
	"os"

	"github.com/mworzala/mc-cli/internal/pkg/app/model"
	"gopkg.in/yaml.v3"
)

func (a *App) OutputFormat() model.OutputFormat {
	format, err := a.cmd.PersistentFlags().GetString("output")
	if err != nil {
		exitWithError(err)
	}

	// Validated by root command pre run function
	return model.OutputFormat(format)
}

func (a *App) LogTyped(m model.PrettyPrintable) {
	var err error
	switch a.OutputFormat() {
	case model.OutputJson:
		err = json.NewEncoder(os.Stdout).Encode(m)
	case model.OutputYaml:
		err = yaml.NewEncoder(os.Stdout).Encode(m)
	default:
		m.Print()
	}
	if err != nil {
		exitWithError(err)
	}
}
