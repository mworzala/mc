//go:build linux
// +build linux

package platform

const ClasspathSeparator = ":"

func GetVersion() (string, error) {
	return "unknown", nil
}
