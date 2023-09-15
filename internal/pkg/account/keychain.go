package account

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"time"

	"github.com/zalando/go-keyring"
)

var ErrNoCredentials = errors.New("no credentials found")

type Credentials struct {
	AccessToken     string    `json:"accessToken"`
	TokenExpiration time.Time `json:"accessExpiration"`
	RefreshToken    string    `json:"refreshToken"`
}

// Keychain is an interface for storing and retrieving account credentials, with support
// for using both a high security system keyring (eg Keychain on macOS, gnome keyring)
// and a local file.
type Keychain interface {
	Get(key string) (*Credentials, error)
	Set(key string, credentials *Credentials) error
}

func NewKeychain(dataDir string, useSystemKeyring bool) Keychain {
	if useSystemKeyring {
		// Try to get a dummy value to test whether the current platform is supported, fall back to file if not
		_, err := keyring.Get(keyringServiceName, "dummy")
		if errors.Is(err, keyring.ErrNotFound) {
			return &systemKeychain{}
		} else if !errors.Is(err, keyring.ErrUnsupportedPlatform) {
			// Some other error
			panic(fmt.Errorf("failed to create system keyring: %w", err))
		}
	}

	// Otherwise use a file, including a warning message once if the file does not exist
	keychainFilePath := path.Join(dataDir, ".keychain")
	if _, err := os.Stat(keychainFilePath); errors.Is(err, fs.ErrNotExist) {
		_, _ = fmt.Fprintln(os.Stderr, keyringUnusedWarning)
	}

	return &fileKeychain{path: keychainFilePath}
}

type fileKeychain struct {
	path string
}

func (k *fileKeychain) Get(key string) (*Credentials, error) {
	panic("not implemented")
}

func (k *fileKeychain) Set(key string, credentials *Credentials) error {
	panic("not implemented")
}

// systemKeychain uses the system keyring to store the credentials.
type systemKeychain struct{}

const (
	keyringServiceName   = "mc-cli.mattworzala.com"
	keyringUnusedWarning = "Warning: account credentials will be stored in plaintext, it is recommended to use a system keyring."
)

func (k *systemKeychain) Get(key string) (*Credentials, error) {
	val, err := keyring.Get(keyringServiceName, key)
	if err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			return nil, ErrNoCredentials
		}
		return nil, fmt.Errorf("failed to get keyring value: %w", err)
	}

	var creds Credentials
	if err = json.Unmarshal([]byte(val), &creds); err != nil {
		return nil, fmt.Errorf("invalid keyring valuer: %w", err)
	}

	return &creds, nil
}

func (k *systemKeychain) Set(key string, credentials *Credentials) error {
	value, err := json.Marshal(credentials)
	if err != nil {
		return fmt.Errorf("failed to marshal credentials: %w", err)
	}

	err = keyring.Set(keyringServiceName, key, string(value))
	if err != nil {
		return fmt.Errorf("failed to set keyring value: %w", err)
	}

	return nil
}
