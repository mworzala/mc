package profile

import (
	"errors"
	"path"

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

	Type Type `json:"type"`
	// Version represents the Minecraft version of the profile.
	// Present no matter the type (except Unknown)
	Version string `json:"version"`
}

func (p *Profile) Config() *Config {
	if p.config != nil {
		return p.config
	}

	v := viper.New()
	v.SetConfigFile(path.Join(p.Directory, "config.json"))

	err := v.ReadInConfig()
	if err != nil && !errors.As(err, &viper.ConfigFileNotFoundError{}) {
		// For some reason viper doesnt implement Is() for ConfigFileNotFoundError
		panic(err) //todo
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		panic(err) //todo
	}

	p.config = &config
	return p.config
}
