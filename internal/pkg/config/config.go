package config

// Config represents the app-global configuration options.
//
// Not all are loaded from the config file, some are set via global command flags.
// Flag values are indicated by a `mapstructure:"-"` tag.
type Config struct {
	NonInteractive bool `mapstructure:"-"`
	//todo output format

	//NoColor      bool             `mapstructure:"no_color"` //todo
	UseSystemKeyring bool             `mapstructure:"use_system_keyring"`
	Experimental     ExperimentalOpts `mapstructure:"experimental"`
}

type ExperimentalOpts struct {
	//todo
}
