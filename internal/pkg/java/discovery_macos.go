//go:build darwin

package java

import (
	"os"
	"path"
)

var javaExecSubPath = "Contents/Home/bin/java"

func discoverKnownPaths() (result []*Installation) {

	// $HOME/Library/Java/JavaVirtualMachines
	if homeDir, err := os.UserHomeDir(); err == nil {
		userLibraryInstalls := path.Join(homeDir, "Library/Java/JavaVirtualMachines")
		for _, install := range discoverDirectory(userLibraryInstalls) {
			result = append(result, install)
		}
	}

	//todo the rest of this
	return
}
