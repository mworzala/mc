package platform

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"runtime"

	"github.com/atotto/clipboard"
)

// GetConfigDir returns the config directory for the cli, or an error
// if the config directory cannot be found.
//
// If the directory does not exist, it will be created.
func GetConfigDir(indev bool) (string, error) {
	var err error
	var configDir string
	if indev {
		configDir, err = os.Getwd()
		if err != nil {
			return "", err
		}
	} else {
		configDir, err = os.UserConfigDir()
		if err != nil {
			return "", err
		}
	}

	dataDir := path.Join(configDir, "mc-cli")

	if _, err := os.Stat(dataDir); errors.Is(err, fs.ErrNotExist) {
		err = os.MkdirAll(dataDir, 0755)
		if err != nil {
			return "", fmt.Errorf("unable to create data directory: %w", err)
		}
	}

	return dataDir, nil
}

// OpenUrl opens the given URL in the system default browser, or returns an error.
// Credit: https://gist.github.com/hyg/9c4afcd91fe24316cbf0
func OpenUrl(url string) error {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", url).Start()
	case "linux":
		return exec.Command("xdg-open", url).Start()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	}
	return fmt.Errorf("cannot open url on %s", runtime.GOOS)
}

func WriteToClipboard(text string) error {
	return clipboard.WriteAll(text)
}
