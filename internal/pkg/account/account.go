package account

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"
	"time"

	"github.com/mworzala/mc-cli/internal/pkg/account/auth"
	"github.com/mworzala/mc-cli/internal/pkg/util"
)

var (
	ErrAccountNotFound = errors.New("no such account")
	accountsFileName   = "accounts.json"
)

type (
	MSOPromptCallback func(verificationUrl, userCode string)
	MSOBeginPolling   func()
)

// todo create a cloudflare worker websocket which can act as a proxy for handling a typical microsoft login,
// then for devicecode require setting `--device-code`. In this case it should check for a user agent on the
// client side. Dont want random people spamming :)

const (
	ModeUUID = 1 << iota
	ModeName = 1 << iota
)

type Manager interface {
	// GetDefault returns the default account, or the empty string if there is no account set
	GetDefault() string
	// SetDefault replaces the default account with the given account, or an error if the account
	// does not exist (ErrAccountNotFound) or another error occurred
	//
	// The Manager is not resposible for persisting the change, Save should be called afterwards.
	SetDefault(uuid string) error

	// GetAccount returns the account with the given value, or nil if it cannot be found.
	// The mode indicates whether to search by ModeUUID, ModeName, or both.
	// When matching on uuid, the form with or without dashes is OK
	// When matching name, it is case insensitive
	GetAccount(value string, mode int) *Account

	// Login mechanisms

	// LoginMicrosoft handles logging into a new Microsoft account
	//
	// The Manager is not responsible for persisting the account to its storage mechanism, Save should be called.
	LoginMicrosoft(promptCallback MSOPromptCallback) (*Account, error)

	Save() error
}

type fileManager struct {
	Path     string              `json:"-"`
	Default  string              `json:"default"`
	Accounts map[string]*Account `json:"accounts"`
}

func NewManager(dataDir string) (Manager, error) {
	accountsFile := path.Join(dataDir, accountsFileName)
	if _, err := os.Stat(accountsFile); errors.Is(err, fs.ErrNotExist) {
		return &fileManager{
			Path:     accountsFile,
			Accounts: make(map[string]*Account),
		}, nil
	}

	f, err := os.Open(accountsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	manager := fileManager{Path: accountsFile}
	if err := json.NewDecoder(f).Decode(&manager); err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", accountsFileName, err)
	}
	if manager.Accounts == nil {
		manager.Accounts = make(map[string]*Account)
	}
	return &manager, nil
}

func (m *fileManager) GetDefault() string {
	return m.Default
}

func (m *fileManager) SetDefault(uuid string) error {
	if _, ok := m.Accounts[uuid]; !ok {
		return ErrAccountNotFound
	}

	m.Default = uuid
	return nil
}

func (m *fileManager) GetAccount(value string, mode int) *Account {
	if mode&ModeUUID != 0 && util.IsUUID(value) {
		return m.Accounts[util.ExpandUUID(value)]
	}

	// If it wasnt a UUID and we are not searching name, we're done
	if mode&ModeName == 0 {
		return nil
	}

	test := strings.ToLower(value)
	for _, acc := range m.Accounts {
		name := strings.ToLower(acc.Profile.Username)
		if test == name {
			return acc
		}
	}
	return nil
}

func (m *fileManager) LoginMicrosoft(promptCallback MSOPromptCallback) (*Account, error) {
	var account Account
	var msoTokenData MicrosoftTokenData
	account.Type = Microsoft
	account.Source = &msoTokenData

	// MSO Device Code
	data, err := auth.BeginDeviceCodeAuth()
	if err != nil {
		return nil, fmt.Errorf("failed to begin device code auth: %w", err)
	}
	promptCallback(data.VerificationURL, data.UserCode)

	msoToken, err := auth.PollDeviceCodeAuth(data)
	if err != nil {
		return nil, fmt.Errorf("failed while polling device code auth: %w", err)
	}
	msoTokenData.AccessToken = msoToken.AccessToken
	msoTokenData.RefreshToken = msoToken.RefreshToken
	msoTokenData.ExpiresAt = time.Now().Add(time.Duration(msoToken.ExpiresIn) * time.Second)

	// Xbox Live
	xblToken, err := auth.XboxLiveAuth(msoToken.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate with Xbox Live: %w", err)
	}
	msoTokenData.UserHash = xblToken.UserHash

	// XSTS
	xstsToken, err := auth.XSTSAuth(xblToken.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate with XSTS: %w", err)
	}

	// Minecraft Auth
	mcToken, err := auth.MinecraftAuthMSO(xstsToken.AccessToken, xblToken.UserHash)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate with Minecraft: %w", err)
	}
	account.AccessToken = mcToken.AccessToken
	account.ExpiresAt = time.Now().Add(time.Duration(mcToken.ExpiresIn) * time.Second)

	// Fetch minecraft profile
	profile, err := auth.GetMinecraftProfile(account.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch minecraft profile: %w", err)
	}
	account.UUID = util.ExpandUUID(profile.UUID)
	account.Profile.Username = profile.Username

	// Make sure to save the newly added account
	m.Accounts[account.UUID] = &account
	return &account, nil
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
