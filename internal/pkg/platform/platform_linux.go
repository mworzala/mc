//go:build linux
// +build linux

package platform

func GetVersion() (string, error) {
	return "unknown", nil
}
