package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Minecraft Auth (mso flow)
// todo doc

type MinecraftToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func MinecraftAuthMSO(xstsToken, userHash string) (*MinecraftToken, error) {
	endpoint := "https://api.minecraftservices.com/authentication/login_with_xbox"
	body := fmt.Sprintf(`{
		"identityToken": "XBL3.0 x=%s;%s",
		"ensureLegacyEnabled": true
	}`, userHash, xstsToken)

	res, err := http.Post(endpoint, "application/json", strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %s", res.Status)
	}

	var payload MinecraftToken
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, err
	}

	return &payload, nil
}
