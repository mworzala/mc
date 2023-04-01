//go:build windows
// +build windows

package platform

import "golang.org/x/sys/windows/registry"

const (
	NTCurrentVersionKey = `SOFTWARE\Microsoft\Windows NT\CurrentVersion`
)

func GetVersion() (string, error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, NTCurrentVersionKey, registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}
	defer k.Close()

	cv, _, err := k.GetStringValue("CurrentVersion")
	return cv, err
}
