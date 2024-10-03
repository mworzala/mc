package skin

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

type Skin struct {
	Name      string    `json:"name"`
	Variant   string    `json:"variant"`
	Skin      string    `json:"skin"`
	Cape      string    `json:"cape"`
	AddedDate time.Time `json:"added_date"`
}

func (s *Skin) Apply(accountToken string) error {
	var newCape bool

	if s.Cape == "none" {
		endpoint := "https://api.minecraftservices.com/minecraft/profile/capes/active"

		req, err := http.NewRequest("DELETE", endpoint, nil)
		if err != nil {
			return err
		}

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accountToken))
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("cape disable request was not ok: %d", res.StatusCode)
		}
	}

	if isURL(s.Skin) {
		endpoint := "https://api.minecraftservices.com/minecraft/profile/skins"

		requestData := map[string]string{
			"url":     s.Skin,
			"variant": s.Variant,
		}

		jsonData, err := json.Marshal(requestData)
		if err != nil {
			return err
		}

		req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
		if err != nil {
			return err
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accountToken))

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("skin url request was not ok: %d", res.StatusCode)
		}

		var information profileInformationResponse
		if err := json.NewDecoder(res.Body).Decode(&information); err != nil {
			return err
		}

		if s.Cape != "none" {
			for _, c := range information.Capes {
				if c.ID == s.Cape && c.State == "INACTIVE" {
					newCape = true
				}
			}
		}

	} else {
		imageData, err := base64.StdEncoding.DecodeString(s.Skin)
		if err != nil {
			return err
		}

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)

		err = writer.WriteField("variant", s.Variant)
		if err != nil {
			return err
		}

		part, err := writer.CreateFormFile("file", "skin.png")
		if err != nil {
			return err
		}

		_, err = io.Copy(part, bytes.NewReader(imageData))
		if err != nil {
			return err
		}

		err = writer.Close()
		if err != nil {
			return err
		}

		req, err := http.NewRequest("POST", "https://api.minecraftservices.com/minecraft/profile/skins", body)
		if err != nil {
			return err
		}

		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accountToken))

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("skin request was not ok: %d", res.StatusCode)
		}

		var information profileInformationResponse
		if err := json.NewDecoder(res.Body).Decode(&information); err != nil {
			return err
		}

		if s.Cape != "none" {
			for _, c := range information.Capes {
				if c.ID == s.Cape && c.State == "INACTIVE" {
					newCape = true
				}
			}
		}

	}

	if newCape {
		endpoint := "https://api.minecraftservices.com/minecraft/profile/capes/active"

		requestData := map[string]string{
			"capeId": s.Cape,
		}

		jsonData, err := json.Marshal(requestData)
		if err != nil {
			return err
		}

		req, err := http.NewRequest("PUT", endpoint, bytes.NewBuffer(jsonData))
		if err != nil {
			return err
		}

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accountToken))
		req.Header.Set("Content-Type", "application/json")
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("cape put request was not ok: %d", res.StatusCode)
		}
	}

	fmt.Printf("skin %s and cape %s was applied", s.Name, s.Cape)

	return nil
}

type profileInformationResponse struct {
	ID             string        `json:"id"`
	Name           string        `json:"name"`
	Skins          []profileSkin `json:"skins"`
	Capes          []profileCape `json:"capes"`
	ProfileActions struct {
	} `json:"profileActions"`
}

type profileSkin struct {
	ID         string `json:"id"`
	State      string `json:"state"`
	URL        string `json:"url"`
	TextureKey string `json:"textureKey"`
	Variant    string `json:"variant"`
}

type profileCape struct {
	ID    string `json:"id"`
	State string `json:"state"`
	URL   string `json:"url"`
	Alias string `json:"alias"`
}
