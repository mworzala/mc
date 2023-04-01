//go:build darwin
// +build darwin

package platform

func GetVersion() (string, error) {
	return "unknown", nil
}
