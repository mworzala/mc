package profile

import (
	"os"
	"path"
	"strings"

	"github.com/spf13/viper"
)

type Type int

const (
	// Unknown indicates that the profile has not been installed yet.
	Unknown Type = iota
	Vanilla
	Fabric
)

type Profile struct {
	Name      string  `json:"name"`
	Directory string  `json:"directory"`
	config    *Config // Config is loaded on demand
	mods      []Mod   `json:"-"` // Stored in Directory/.mc-cli/mods.json

	Type Type `json:"type"`
	// Version represents the version json used for the profile.
	// Present no matter the type (except Unknown)
	Version string `json:"version"`

	// GameVersion represents the Minecraft version of the profile.
	// Accessed through GameVersion() because it defaults to parsing Version in the event that the profile is old
	gameVersion string `json:"game_version"`
}

func (p *Profile) Config() *Config {
	if p.config != nil {
		return p.config
	}

	v := viper.New()
	v.SetConfigFile(path.Join(p.Directory, "config.json"))

	// Read the config file if it exists
	if _, err := os.Stat(v.ConfigFileUsed()); err == nil {
		if err := v.ReadInConfig(); err != nil {
			panic(err) //todo this should proceed with default config and log an error
		}
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		panic(err) //todo
	}

	p.config = &config
	return p.config
}

func (p *Profile) GameVersion() string {
	if p.gameVersion != "" {
		return p.gameVersion
	}

	splits := strings.Split(p.Version, "-")
	if len(splits) < 4 {
		return p.Version
	}

	p.gameVersion = splits[3]
	return p.gameVersion
}
