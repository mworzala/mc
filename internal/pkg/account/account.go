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

	"github.com/mworzala/mc/internal/pkg/config"

	"github.com/mworzala/mc/internal/pkg/account/auth"
	"github.com/mworzala/mc/internal/pkg/util"
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

type Manager interface {
	// GetDefault returns the default account, or the empty string if there is no account set
	GetDefault() string
	// SetDefault replaces the default account with the given account, or an error if the account
	// does not exist (ErrAccountNotFound) or another error occurred
	//
	// The Manager is not resposible for persisting the change, Save should be called afterwards.
	SetDefault(uuid string) error

	Accounts() []string
	// GetAccount returns the account with the given value, or nil if it cannot be found.
	// Either a (case-insensitive) name, or a UUID (with/without dashes) can be matched
	GetAccount(value string) *Account
	// GetAccountToken returns a _minecraft_ access token for the given account.
	// The given value may be a (case-insensitive) name, or a UUID (with/without dashes).
	//
	// This function will always return an active token, using the refresh token if necessary.
	GetAccountToken(value string) (string, error)

	// Login mechanisms

	// LoginMicrosoft handles logging into a new Microsoft account
	//
	// The Manager is not responsible for persisting the account to its storage mechanism, Save should be called.
	LoginMicrosoft(promptCallback MSOPromptCallback) (*Account, error)

	Save() error
}

type fileManager struct {
	Path     string   `json:"-"`
	Keychain Keychain `json:"-"`

	Default     string              `json:"default"`
	AccountData map[string]*Account `json:"accounts"`
}

func NewManager(dataDir string, config *config.Config) (Manager, error) {

	// Read the accounts file
	accountsFile := path.Join(dataDir, accountsFileName)
	if _, err := os.Stat(accountsFile); errors.Is(err, fs.ErrNotExist) {
		return &fileManager{
			Path:        accountsFile,
			AccountData: make(map[string]*Account),
		}, nil
	}
	f, err := os.Open(accountsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	// Construct the manager
	manager := fileManager{Path: accountsFile, AccountData: make(map[string]*Account)}
	if err := json.NewDecoder(f).Decode(&manager); err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", accountsFileName, err)
	}

	manager.Keychain = NewKeychain(dataDir, config.UseSystemKeyring)

	return &manager, nil
}

func (m *fileManager) GetDefault() string {
	return m.Default
}

func (m *fileManager) SetDefault(uuid string) error {
	if _, ok := m.AccountData[uuid]; !ok {
		return ErrAccountNotFound
	}

	m.Default = uuid
	return nil
}

func (m *fileManager) Accounts() (result []string) {
	for k := range m.AccountData {
		result = append(result, k)
	}
	return
}

func (m *fileManager) GetAccount(value string) *Account {
	isUuid := util.IsUUID(value)
	for uuid, account := range m.AccountData {
		if isUuid && uuid == util.ExpandUUID(value) {
			return account
		}

		if !isUuid && strings.ToLower(account.Profile.Username) == strings.ToLower(value) {
			return account
		}
	}

	return nil
}

func (m *fileManager) GetAccountToken(value string) (string, error) {
	account := m.GetAccount(value)
	if account == nil {
		return "", ErrAccountNotFound
	}

	credentials, err := m.Keychain.Get(account.UUID)
	if err != nil {
		return "", fmt.Errorf("failed to get credentials: %w", err)
	}

	// Update the credentials if necessary
	changed, err := m.updateCredentialsMso(credentials)
	if err != nil {
		return "", fmt.Errorf("failed to update credentials: %w", err)
	}
	if changed {
		if err := m.Keychain.Set(account.UUID, credentials); err != nil {
			return "", fmt.Errorf("failed to save credentials: %w", err)
		}
	}

	return credentials.AccessToken, nil
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

	userHash, accessToken, tokenExpiration, err := m.createCredentialsMso(msoToken.AccessToken)
	if err != nil {
		return nil, err
	}
	msoTokenData.UserHash = userHash

	// Fetch minecraft profile
	profile, err := auth.GetMinecraftProfile(accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch minecraft profile: %w", err)
	}
	account.UUID = util.ExpandUUID(profile.UUID)
	account.Profile.Username = profile.Username

	// Make sure to save the newly added account
	err = m.Keychain.Set(account.UUID, &Credentials{
		AccessToken:     accessToken,
		TokenExpiration: tokenExpiration,
		RefreshToken:    msoToken.RefreshToken,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to save account to keychain: %w", err)
	}
	m.AccountData[account.UUID] = &account

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

// updateCredentialsMso attempts to update a set of MSO credentials using the refresh token if the included
// minecraft access token is missing or expired.
func (m *fileManager) updateCredentialsMso(credentials *Credentials) (bool, error) {
	// If the access token is present and not expired, do nothing
	if credentials.AccessToken != "" && credentials.TokenExpiration.After(time.Now()) {
		return false, nil
	}
	//todo would be nice to log a message here about refreshing the token in a debug log

	// Sanity check
	if credentials.RefreshToken == "" {
		return false, errors.New("missing refresh token")
	}

	// Refresh the credentials
	msoToken, err := auth.RefreshMsoToken(credentials.RefreshToken)
	if err != nil {
		return false, fmt.Errorf("failed to refresh MSO token: %w", err)
	}

	credentials.RefreshToken = msoToken.RefreshToken

	// Get a new Minecraft access token
	_, credentials.AccessToken, credentials.TokenExpiration, err = m.createCredentialsMso(msoToken.AccessToken)
	if err != nil {
		return false, err
	}

	return true, nil
}

// createMsoCredentials creates a new minecraft access token from an MSO access token.
func (m *fileManager) createCredentialsMso(msoAccessToken string) (xblUserHash string, accessToken string, tokenExpiration time.Time, err error) {
	// Xbox Live
	xblToken, err := auth.XboxLiveAuth(msoAccessToken)
	if err != nil {
		err = fmt.Errorf("failed to authenticate with Xbox Live: %w", err)
		return
	}
	xblUserHash = xblToken.UserHash

	// XSTS
	xstsToken, err := auth.XSTSAuth(xblToken.AccessToken)
	if err != nil {
		err = fmt.Errorf("failed to authenticate with XSTS: %w", err)
		return
	}

	// Minecraft Auth
	mcToken, err := auth.MinecraftAuthMSO(xstsToken.AccessToken, xblToken.UserHash)
	if err != nil {
		err = fmt.Errorf("failed to authenticate with Minecraft: %w", err)
		return
	}
	accessToken = mcToken.AccessToken
	tokenExpiration = time.Now().Add(time.Duration(mcToken.ExpiresIn) * time.Second)

	return
}
