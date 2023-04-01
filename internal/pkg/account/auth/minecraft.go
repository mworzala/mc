package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Minecraft Auth (mso flow)
// https://mojang-api-docs.gapple.pw/authentication/msa#getting-the-bearer-token-for-minecraft

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

type MinecraftProfile struct {
	UUID     string `json:"id"`
	Username string `json:"name"`
}

func GetMinecraftProfile(accessToken string) (*MinecraftProfile, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api.minecraftservices.com/minecraft/profile", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %w", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer res.Body.Close()

	var profile MinecraftProfile
	if err := json.NewDecoder(res.Body).Decode(&profile); err != nil {
		return nil, fmt.Errorf("failed to parse profile response body: %w", err)
	}

	return &profile, nil
}
