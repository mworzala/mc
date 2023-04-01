package account

import (
	"encoding/json"
	"fmt"
	"time"
)

type Type string

const (
	Microsoft Type = "microsoft"
	Mojang    Type = "mojang"
)

func TypeFromString(accountType string) (Type, bool) {
	switch accountType {
	case "microsoft", "mso":
		return Microsoft, true
	case "mojang", "minecraft", "mc":
		return Mojang, true
	}
	return "", false
}

type Account struct {
	UUID    string `json:"uuid"`
	Profile struct {
		Username string `json:"username"`
	} `json:"profile"`

	// Minecraft auth token, the source of the token is Token
	AccessToken string `json:"accessToken"`
	// Expiry time for AccessToken
	ExpiresAt time.Time `json:"expiresAt"`

	Type Type `json:"type"`
	// Source is the token data for the auth Type.
	// Either MicrosoftTokenData or (todo) MojangTokenData.
	Source interface{} `json:"source"`
}

func (a *Account) UnmarshalJSON(data []byte) error {
	type Delegate Account
	type AccountWithRawSource struct {
		Delegate
		SourceRaw json.RawMessage `json:"source"`
	}
	var acc AccountWithRawSource
	if err := json.Unmarshal(data, &acc); err != nil {
		return err
	}

	if acc.Type == Microsoft {
		var msoTokenData MicrosoftTokenData
		if err := json.Unmarshal(acc.SourceRaw, &msoTokenData); err != nil {
			return fmt.Errorf("failed to unmarshal microsoft token: %w", err)
		}
		acc.Source = &msoTokenData
	} else {
		panic("not implemented")
	}

	*a = Account(acc.Delegate)
	return nil
}

type MicrosoftTokenData struct {
	UserHash     string    `json:"userHash"`
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	ExpiresAt    time.Time `json:"expiresAt"`
}
