package skin

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

var (
	ErrInvalidName = errors.New("invalid skin name")
	ErrNameInUse   = errors.New("name in use")
	ErrNotFound    = errors.New("skin not found")

	namePattern = regexp.MustCompile("^[a-zA-Z0-9_.-]{1,36}$")
)

func isValidName(name string) bool {
	return namePattern.MatchString(name)
}

func isURL(s string) bool {
	u, err := url.ParseRequestURI(s)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https")
}

func isFilePath(s string) bool {
	if _, err := os.Stat(s); err == nil {
		return true
	}
	return false
}

func isImage(data []byte) bool {
	contentType := http.DetectContentType(data)

	// probably more of these however i know png and jpeg are supported
	if contentType == "image/png" || contentType == "image/jpeg" {
		return true
	}

	return false
}

type Manager interface {
	CreateSkin(name string, variant string, skinData string, capeData string) (*Skin, error)
	Skins() []*Skin
	GetSkin(name string) (*Skin, error)

	GetProfileInformation(accountToken string) (*profileInformationResponse, error)

	Save() error
}

var (
	skinsFileName = "skins.json"
)

type fileManager struct {
	Path     string           `json:"-"`
	AllSkins map[string]*Skin `json:"skins"`
}

func NewManager(dataDir string) (Manager, error) {
	skinsFile := path.Join(dataDir, skinsFileName)
	if _, err := os.Stat(skinsFile); errors.Is(err, fs.ErrNotExist) {
		return &fileManager{
			Path:     skinsFile,
			AllSkins: make(map[string]*Skin),
		}, nil
	}

	f, err := os.Open(skinsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	manager := fileManager{Path: skinsFile}
	if err := json.NewDecoder(f).Decode(&manager); err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", skinsFileName, err)
	}
	if manager.AllSkins == nil {
		manager.AllSkins = make(map[string]*Skin)
	}
	return &manager, nil
}

func (m *fileManager) CreateSkin(name string, variant string, skinData string, capeData string) (*Skin, error) {
	if !isValidName(name) {
		return nil, ErrInvalidName
	}
	if _, ok := m.AllSkins[strings.ToLower(name)]; ok {
		return nil, ErrNameInUse
	}

	skin := &Skin{
		Name:    name,
		Cape:    capeData,
		Variant: variant,
	}
	if isFilePath(skinData) {
		fileBytes, err := os.ReadFile(skinData)
		if err != nil {
			return nil, err
		}
		isValid := isImage(fileBytes)
		if !isValid {
			return nil, fmt.Errorf("%s is not a valid image", skinData)
		}

		base64Str := base64.StdEncoding.EncodeToString(fileBytes)
		skin.Skin = base64Str
	} else {
		skin.Skin = skinData
	}

	skin.AddedDate = time.Now()

	m.AllSkins[strings.ToLower(name)] = skin
	return skin, nil
}

func (m *fileManager) Skins() (result []*Skin) {
	for _, s := range m.AllSkins {
		result = append(result, s)
	}
	return
}

func (m *fileManager) GetSkin(name string) (*Skin, error) {
	name = strings.ToLower(name)
	var matchingSkin *Skin
	matchCount := 0
	for id, s := range m.AllSkins {
		if id == name {
			return s, nil
		}
		if strings.HasPrefix(id, name) {
			matchCount++
			matchingSkin = s
		}
	}

	if matchCount == 1 {
		return matchingSkin, nil
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

func (m *fileManager) GetProfileInformation(accountToken string) (*profileInformationResponse, error) {
	endpoint := "https://api.minecraftservices.com/minecraft/profile"

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accountToken))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var information profileInformationResponse
	if err := json.NewDecoder(res.Body).Decode(&information); err != nil {
		return nil, err
	}

	return &information, nil

}
