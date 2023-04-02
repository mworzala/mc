package java

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
)

func discoverDirectory(dir string) []*Installation {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	var result []*Installation
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		execPath := path.Join(dir, entry.Name(), javaExecSubPath)
		if _, err := os.Stat(execPath); errors.Is(err, fs.ErrNotExist) {
			continue
		}

		// Error doesn't matter, just move on
		install, _ := DiscoverJava(execPath, entry.Name())
		result = append(result, install)
	}

	return result
}

func DiscoverJava(executable, name string) (*Installation, error) {
	var i Installation
	i.Path = executable
	i.Name = name

	if err := discoverParams(&i); err != nil {
		return nil, err
	}

	return &i, nil
}

func discoverParams(i *Installation) error {
	cmd := exec.Command(i.Path, "-XshowSettings:properties", "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("java discovery exec error: %w", err)
	}

	properties := extractParams(string(output))
	i.Arch = properties["os.arch"]
	v, err := strconv.ParseInt(properties["java.vm.specification.version"], 10, 64)
	if err != nil {
		return err
	}
	i.Version = int(v)

	if i.Name == "" {
		i.Name = properties["java.vendor.version"]
	}

	return nil
}

var paramRegex = regexp.MustCompile(`(?m)^([^\n=]+)\s*=\s*(.*)$`)

func extractParams(text string) map[string]string {
	result := make(map[string]string)
	matches := paramRegex.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		result[strings.TrimSpace(match[1])] = strings.TrimSpace(match[2])
	}
	return result
}
