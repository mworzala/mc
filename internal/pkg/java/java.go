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

	Discover(name string) (*Installation, error)

	Installations() []string
	GetInstallation(name string) *Installation

	Save() error
}

type fileManager struct {
	Path     string                   `json:"-"`
	Default  string                   `json:"default"`
	Installs map[string]*Installation `json:"installs"`
}

func NewManager(dataDir string) (Manager, error) {
	javaFile := path.Join(dataDir, installsFileName)
	if _, err := os.Stat(javaFile); errors.Is(err, fs.ErrNotExist) {
		installs := make(map[string]*Installation)

		for _, i := range discoverKnownPaths() {
			installs[strings.ToLower(i.Name)] = i
		}

		return &fileManager{
			Path:     javaFile,
			Installs: installs,
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
	if manager.Installs == nil {
		manager.Installs = make(map[string]*Installation)
	}
	return &manager, nil
}

func (m *fileManager) GetDefault() string {
	return m.Default
}

func (m *fileManager) SetDefault(name string) error {
	name = strings.ToLower(name)
	if _, ok := m.Installs[name]; !ok {
		return ErrInstallationNotFound
	}

	m.Default = name
	return nil
}

func (m *fileManager) Discover(exec string) (*Installation, error) {
	install, err := DiscoverJava(exec, "")
	if err != nil {
		return nil, err
	}

	m.Installs[strings.ToLower(install.Name)] = install
	return install, err
}

func (m *fileManager) Installations() (result []string) {
	for k := range m.Installs {
		result = append(result, k)
	}
	return
}

func (m *fileManager) GetInstallation(name string) *Installation {
	return m.Installs[strings.ToLower(name)]
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
