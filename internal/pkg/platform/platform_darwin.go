//go:build darwin
// +build darwin

package platform

const (
	Name               = "osx"
	ClasspathSeparator = ":"
)

func GetVersion() (string, error) {
	return "unknown", nil
}
