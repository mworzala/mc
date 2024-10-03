package modrinth

type ProjectType string

const (
	Mod          ProjectType = "mod"
	ModPack      ProjectType = "modpack"
	ResourcePack ProjectType = "resourcepack"
	Shader       ProjectType = "shader"
)

type SupportStatus string

const (
	Required    SupportStatus = "required"
	Optional    SupportStatus = "optional"
	Unsupported SupportStatus = "unsupported"
	Unknown     SupportStatus = "unknown"
)

type MonetizationStatus string

const (
	Monetized        MonetizationStatus = "monetized"
	Demonetized      MonetizationStatus = "demonetized"
	ForceDemonetized MonetizationStatus = "force-demonetized"
)
