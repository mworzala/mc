//go:build darwin
// +build darwin

package platform

const ClasspathSeparator = ":"

func GetVersion() (string, error) {
	return "unknown", nil
}
