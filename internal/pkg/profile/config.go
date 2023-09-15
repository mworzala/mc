package profile

// Config represents all the profile specific configuration options.
// Options not specified in a profile config will be inherited from the global config.
type Config struct {
	Java string `mapstructure:"java"` // The name of the java installation to use
}
