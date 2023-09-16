//go:build linux
// +build linux

package platform

const (
	Name               = "linux"
	ClasspathSeparator = ":"
)

func GetVersion() (string, error) {
	return "unknown", nil
}
