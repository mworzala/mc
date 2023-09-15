package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	msoTenant    = "consumers"
	msoClientId  = "5b71e97c-5415-4e16-8bdb-9c638a939986"
	msoGrantType = "urn:ietf:params:oauth:grant-type:device_code"
	msoScope     = "XboxLive.signin offline_access"
)

type MsoTokenData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// Microsoft OAuth (device code access flow)
// https://docs.microsoft.com/en-us/azure/active-directory/develop/v2-oauth2-device-code

type DeviceCodeData struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURL string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

type deviceCodePollErrorResponse struct {
	Error       deviceCodePollError `json:"error"`
	Description string              `json:"error_description"`
}

type deviceCodePollError string

const (
	errorAuthorizationPending deviceCodePollError = "authorization_pending"
)

func BeginDeviceCodeAuth() (*DeviceCodeData, error) {
	endpoint := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/devicecode", msoTenant)
	body := fmt.Sprintf("client_id=%s&scope=%s", url.QueryEscape(msoClientId), url.QueryEscape(msoScope))

	res, err := http.Post(endpoint, "application/x-www-form-urlencoded", strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %s", res.Status)
	}

	var payload DeviceCodeData
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, err
	}
	return &payload, nil
}

func PollDeviceCodeAuth(data *DeviceCodeData) (*MsoTokenData, error) {
	endpoint := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", msoTenant)
	body := fmt.Sprintf("grant_type=%s&client_id=%s&device_code=%s",
		url.QueryEscape(msoGrantType), url.QueryEscape(msoClientId), url.QueryEscape(data.DeviceCode))

	expiry := time.Now().Add(time.Duration(data.ExpiresIn) * time.Second)
	//todo probably could use a timeout context to handle the expiration here (more accurately)
	for time.Now().Before(expiry) {
		res, err := http.Post(endpoint, "application/x-www-form-urlencoded", strings.NewReader(body))
		if err != nil {
			return nil, err
		}
		//goland:noinspection GoDeferInLoop
		defer res.Body.Close()

		// If OK, we're done
		switch res.StatusCode {
		case http.StatusOK:
			var payload MsoTokenData
			if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
				return nil, err
			}

			return &payload, nil
		case http.StatusBadRequest:
			var payload deviceCodePollErrorResponse
			if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
				return nil, err
			}

			if payload.Error == errorAuthorizationPending {
				time.Sleep(time.Duration(data.Interval) * time.Second)
				continue
			}

			return nil, fmt.Errorf("device code auth error: %s", payload.Description)
		default:
			return nil, fmt.Errorf("HTTP error: %s", res.Status)
		}
	}

	return nil, fmt.Errorf("device code expired") //todo define as constant
}

// Microsoft OAuth (refresh token flow)
// https://learn.microsoft.com/en-us/azure/active-directory/develop/v2-oauth2-auth-code-flow#refresh-the-access-token

func RefreshMsoToken(refreshToken string) (*MsoTokenData, error) {
	endpoint := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", msoTenant)
	body := fmt.Sprintf("client_id=%s&scope=%s&refresh_token=%s&grant_type=refresh_token",
		url.QueryEscape(msoClientId), url.QueryEscape(msoScope), url.QueryEscape(refreshToken))

	res, err := http.Post(endpoint, "application/x-www-form-urlencoded", strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %s", res.Status)
	}

	var payload MsoTokenData
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, err
	}
	return &payload, nil
}

// Xbox Live (mso to xbl token exchange)
// https://mojang-api-docs.netlify.app/authentication/msa.html#signing-into-xbox-live

type XboxLiveToken struct {
	AccessToken string `json:"access_token"`
	UserHash    string `json:"user_hash"`
}

type xboxLiveAuthResponse struct {
	DisplayClaims struct {
		Xui []struct {
			Uhs string `json:"uhs"`
		} `json:"xui"`
	} `json:"DisplayClaims"`
	Token string `json:"Token"`
}

func XboxLiveAuth(msoAccessToken string) (*XboxLiveToken, error) {
	endpoint := "https://user.auth.xboxlive.com/user/authenticate"
	body := fmt.Sprintf(`{
"Properties": {
"AuthMethod": "RPS",
"SiteName": "user.auth.xboxlive.com",
"RpsTicket": "d=%s"
},
"RelyingParty": "http://auth.xboxlive.com",
"TokenType": "JWT"
}`, msoAccessToken)

	res, err := http.Post(endpoint, "application/json", strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %s", res.Status)
	}

	var payload xboxLiveAuthResponse
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, err
	}

	return &XboxLiveToken{
		AccessToken: payload.Token,
		UserHash:    payload.DisplayClaims.Xui[0].Uhs,
	}, nil
}

// XSTS (xbl to xsts token exchange)
// https://mojang-api-docs.netlify.app/authentication/msa.html#getting-an-xsts-token

type XSTSToken struct {
	AccessToken string `json:"Token"`
}

func XSTSAuth(xblToken string) (*XSTSToken, error) {
	endpoint := "https://xsts.auth.xboxlive.com/xsts/authorize"
	body := fmt.Sprintf(`{
"Properties": {
"SandboxId": "RETAIL",
"UserTokens": ["%s"]
},
"RelyingParty": "rp://api.minecraftservices.com/",
"TokenType": "JWT"
}`, xblToken)

	res, err := http.Post(endpoint, "application/json", strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %s", res.Status)
	}

	var payload XSTSToken
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, err
	}

	return &payload, nil
}
