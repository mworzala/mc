package profile

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"regexp"
	"strings"
)

var (
	ErrInvalidName = errors.New("invalid profile name")
	ErrNameInUse   = errors.New("name in use")
	ErrNotFound    = errors.New("profile not found")

	namePattern = regexp.MustCompile("^[a-zA-Z0-9_.-]{1,32}$")
)

func IsValidName(name string) bool {
	//todo tests
	return namePattern.MatchString(name)
}

type Manager interface {
	// CreateProfile creates a new profile and fills in defaults.
	// The returned profile may be modified, and then Save will save it.
	CreateProfile(name string) (*Profile, error)

	Profiles() []string
	GetProfile(name string) (*Profile, error)

	Save() error
}

var (
	profilesFileName = "profiles.json"
)

type fileManager struct {
	Path        string `json:"-"`
	profilesDir string
	Default     string              `json:"default"`
	AllProfiles map[string]*Profile `json:"profiles"`
}

func NewManager(dataDir string) (Manager, error) {
	profilesFile := path.Join(dataDir, profilesFileName)
	if _, err := os.Stat(profilesFile); errors.Is(err, fs.ErrNotExist) {
		return &fileManager{
			Path:        profilesFile,
			profilesDir: path.Join(dataDir, "profiles"),
			AllProfiles: make(map[string]*Profile),
		}, nil
	}

	f, err := os.Open(profilesFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open profiles file: %w", err)
	}
	defer f.Close()

	manager := fileManager{Path: profilesFile, profilesDir: path.Join(dataDir, "profiles")}
	if err := json.NewDecoder(f).Decode(&manager); err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", profilesFileName, err)
	}
	if manager.AllProfiles == nil {
		manager.AllProfiles = make(map[string]*Profile)
	}
	return &manager, nil
}

func (m *fileManager) CreateProfile(name string) (*Profile, error) {
	if !IsValidName(name) {
		return nil, ErrInvalidName
	}
	if _, ok := m.AllProfiles[strings.ToLower(name)]; ok {
		return nil, ErrNameInUse
	}

	dataDir := path.Join(m.profilesDir, name)
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create profile data directory: %w", err)
	}

	profile := &Profile{
		Name:      name,
		Type:      Unknown,
		Directory: dataDir,
	}

	m.AllProfiles[strings.ToLower(name)] = profile
	return profile, nil
}

func (m *fileManager) Profiles() (result []string) {
	for name := range m.AllProfiles {
		result = append(result, name)
	}
	return
}

func (m *fileManager) GetProfile(name string) (*Profile, error) {
	name = strings.ToLower(name)
	for id, p := range m.AllProfiles {
		if id == name {
			return p, nil
		}
	}
	return nil, ErrNotFound
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
