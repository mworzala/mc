package account

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type Type string

const (
	Microsoft Type = "microsoft"
	Mojang    Type = "mojang"
)

var ErrInvalidType = errors.New("invalid account type")

func ParseType(s string) (Type, error) {
	switch strings.ToLower(s) {
	case "microsoft", "mso":
		return Microsoft, nil
	case "mojang", "minecraft", "mc":
		return Mojang, nil
	}
	return "", ErrInvalidType
}

type Account struct {
	UUID    string `json:"uuid"`
	Profile struct {
		Username string `json:"username"`
	} `json:"profile"`

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
	UserHash string `json:"userHash"`
}
