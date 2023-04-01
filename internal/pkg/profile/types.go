package profile

type Type int

const (
	// Unknown indicates that the profile has not been installed yet.
	Unknown Type = iota
	Vanilla
	Fabric
)

type Profile struct {
	Name      string `json:"name"`
	Directory string `json:"directory"`

	Type Type `json:"type"`
	// Version represents the Minecraft version of the profile.
	// Present no matter the type (except Unknown)
	Version string `json:"version"`
}
