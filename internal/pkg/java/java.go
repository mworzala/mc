package java

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"
)

var (
	ErrInstallationNotFound = errors.New("installation not found")
	installsFileName        = "java.json"
)

type Installation struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Arch    string `json:"arch"`
	Version int    `json:"version"`
}

type Manager interface {
	GetDefault() string
	SetDefault(name string) error

	GetInstallation(name string) *Installation

	Save() error
}

type fileManager struct {
	Path          string                   `json:"-"`
	Default       string                   `json:"default"`
	Installations map[string]*Installation `json:"installs"`
}

func NewManager(dataDir string) (Manager, error) {
	javaFile := path.Join(dataDir, installsFileName)
	if _, err := os.Stat(javaFile); errors.Is(err, fs.ErrNotExist) {
		installs := make(map[string]*Installation)

		for _, i := range DefaultDiscoverMacOS() {
			println("Discovered", i.Path, i.Arch)
			installs[strings.ToLower(i.Name)] = i
		}

		return &fileManager{
			Path:          javaFile,
			Installations: installs,
		}, nil
	}

	f, err := os.Open(javaFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	manager := fileManager{Path: javaFile}
	if err := json.NewDecoder(f).Decode(&manager); err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", installsFileName, err)
	}
	if manager.Installations == nil {
		manager.Installations = make(map[string]*Installation)
	}
	return &manager, nil
}

func (m *fileManager) GetDefault() string {
	return m.Default
}

func (m *fileManager) SetDefault(name string) error {
	name = strings.ToLower(name)
	if _, ok := m.Installations[name]; !ok {
		return ErrInstallationNotFound
	}

	m.Default = name
	return nil
}

func (m *fileManager) GetInstallation(name string) *Installation {
	return m.Installations[strings.ToLower(name)]
}

func (m *fileManager) Save() error {
	f, err := os.OpenFile(m.Path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("failed to open %s: %w", m.Path, err)
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(m); err != nil {
		return fmt.Errorf("failed to write json: %w", err)
	}

	return nil
}
