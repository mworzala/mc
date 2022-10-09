package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Provider struct {
	rootDir string
	loaded  bool
	// accountInfo maps uuid to account info
	accountInfo map[string]msoAccountInfo
}

func NewProvider(rootDir string) Provider {
	return Provider{
		rootDir:     rootDir,
		loaded:      false,
		accountInfo: map[string]msoAccountInfo{},
	}
}

// msoAccountInfo is a stored account info for a MSO account
type msoAccountInfo struct {
	UUID string `json:"uuid"`
	//todo Profile

	// MSO token info
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	ExpiresAt    time.Time `json:"expiresAt"`

	// Latest is the latest minecraft token, used if not exired
	Latest mcTokenInfo `json:"latest"`
}

type mcTokenInfo struct {
	AccessToken string    `json:"accessToken"`
	ExpiresAt   time.Time `json:"expiresAt"`
	UserHash    string    `json:"userHash"`
}

type Credentials struct {
	PlayerName  string
	UUID        string
	UserType    UserType
	AccessToken string
	UserHash    string
}

type UserType string

var (
	UserTypeMSA UserType = "msa"
)

var ErrNotFound = errors.New("account not found")

func (p *Provider) GetCredentials(uuid string) (Credentials, error) {
	if err := p.loadAccountInfo(); err != nil {
		return Credentials{}, fmt.Errorf("failed to load account info: %w", err)
	}

	account, ok := p.accountInfo[uuid]
	if !ok {
		return Credentials{}, ErrNotFound
	}

	// Check if latest token is expired
	if account.Latest.ExpiresAt.Before(time.Now()) {
		//todo refresh the token
		return Credentials{}, errors.New("token expired")
	}

	return Credentials{
		PlayerName:  "notmattw",
		UUID:        uuid,
		UserType:    UserTypeMSA,
		AccessToken: account.Latest.AccessToken,
		UserHash:    account.Latest.UserHash,
	}, nil
}

func (p *Provider) LoginMSA() (string, error) {
	var account msoAccountInfo

	// MSO Device Code
	data, err := BeginDeviceCodeAuth()
	if err != nil {
		return "", fmt.Errorf("failed to begin device code auth: %w", err)
	}
	fmt.Println(data.VerificationURL + " " + data.UserCode)

	msoToken, err := PollDeviceCodeAuth(data)
	if err != nil {
		return "", fmt.Errorf("failed while polling device code auth: %w", err)
	}
	account.AccessToken = msoToken.AccessToken
	account.RefreshToken = msoToken.RefreshToken
	account.ExpiresAt = time.Now().Add(time.Duration(msoToken.ExpiresIn) * time.Second)

	// Xbox Live
	xblToken, err := XboxLiveAuth(msoToken.AccessToken)
	if err != nil {
		return "", fmt.Errorf("failed to authenticate with Xbox Live: %w", err)
	}
	account.Latest.UserHash = xblToken.UserHash

	// XSTS
	xstsToken, err := XSTSAuth(xblToken.AccessToken)
	if err != nil {
		return "", fmt.Errorf("failed to authenticate with XSTS: %w", err)
	}

	// Minecraft Auth
	mcToken, err := MinecraftAuthMSO(xstsToken.AccessToken, xblToken.UserHash)
	if err != nil {
		return "", fmt.Errorf("failed to authenticate with Minecraft: %w", err)
	}
	account.Latest.AccessToken = mcToken.AccessToken
	account.Latest.ExpiresAt = time.Now().Add(time.Duration(mcToken.ExpiresIn) * time.Second)

	// Fetch minecraft profile info todo
	account.UUID = "aceb326fda1545bcbf2f11940c21780c"

	// Save to storage
	if err = p.loadAccountInfo(); err != nil {
		return "", fmt.Errorf("failed to load account info: %w", err)
	}
	p.accountInfo[account.UUID] = account
	if err = p.saveAccountInfo(); err != nil {
		return "", fmt.Errorf("failed to save account info: %w", err)
	}

	return "notmattw", nil
}

// Helpers

func (p *Provider) loadAccountInfo() error {
	if p.loaded {
		return nil
	}

	// Load from storage
	authFile := filepath.Join(p.rootDir, "auth.json")
	f, err := os.OpenFile(authFile, os.O_RDONLY, 0600)
	if err != nil {
		if os.IsNotExist(err) {
			p.loaded = true
			return nil
		}
		return err
	}

	// Parse json
	var accounts map[string]msoAccountInfo
	if err := json.NewDecoder(f).Decode(&accounts); err != nil {
		return err
	}

	p.loaded = true
	p.accountInfo = accounts
	return nil
}

func (p *Provider) saveAccountInfo() error {
	authFile := filepath.Join(p.rootDir, "auth.json")
	f, err := os.OpenFile(authFile, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	// Encode json
	if err := json.NewEncoder(f).Encode(p.accountInfo); err != nil {
		return err
	}

	return nil
}
