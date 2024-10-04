package mojang

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
)

func isURL(s string) bool {
	u, err := url.ParseRequestURI(s)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https")
}

func (c *Client) ProfileInformation(ctx context.Context, accountToken string) (*ProfileInformationResponse, error) {

	headers := http.Header{}

	headers.Set("Authorization", "Bearer "+accountToken)

	return get[ProfileInformationResponse](c, ctx, "", headers)
}

func (c *Client) ChangeSkin(ctx context.Context, accountToken string, texture string, variant string) (*ProfileInformationResponse, error) {
	var body *bytes.Buffer
	var contentType string

	if isURL(texture) {
		requestData := map[string]string{
			"url":     texture,
			"variant": variant,
		}

		jsonData, err := json.Marshal(requestData)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(jsonData)
		contentType = "application/json"
	} else {
		imageData, err := base64.StdEncoding.DecodeString(texture)
		if err != nil {
			return nil, err
		}

		body = new(bytes.Buffer)
		writer := multipart.NewWriter(body)

		err = writer.WriteField("variant", variant)
		if err != nil {
			return nil, err
		}

		part, err := writer.CreateFormFile("file", "skin.png")
		if err != nil {
			return nil, err
		}

		_, err = io.Copy(part, bytes.NewReader(imageData))
		if err != nil {
			return nil, err
		}

		err = writer.Close()
		if err != nil {
			return nil, err
		}

		contentType = writer.FormDataContentType()
	}

	headers := http.Header{}

	headers.Set("Content-Type", contentType)
	headers.Set("Authorization", "Bearer "+accountToken)

	return post[ProfileInformationResponse](c, ctx, "/skins", headers, body)
}

func (c *Client) ChangeCape(ctx context.Context, accountToken string, cape string) (*ProfileInformationResponse, error) {
	endpoint := "capes/active"
	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+accountToken)
	headers.Set("Content-Type", "application/json")

	requestData := map[string]string{
		"capeId": cape,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return nil, err
	}

	return put[ProfileInformationResponse](c, ctx, endpoint, headers, bytes.NewBuffer(jsonData))
}

func (c *Client) DeleteCape(ctx context.Context, accountToken string) (*ProfileInformationResponse, error) {
	endpoint := "capes/active"
	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+accountToken)

	return delete[ProfileInformationResponse](c, ctx, endpoint, headers)
}

func (c *Client) UsernameToUuid(ctx context.Context, username string) (*UsernameToUuidResponse, error) {
	oldUrl := c.baseUrl
	c.baseUrl = mojangApiUrl // i dont like this but i cant think of any other way atm :(
	endpoint := "users/profiles/minecraft/" + username

	response, err := get[UsernameToUuidResponse](c, ctx, endpoint, http.Header{})
	c.baseUrl = oldUrl
	return response, err
}

func (c *Client) UuidToProfile(ctx context.Context, uuid string) (*UuidToProfileResponse, error) {
	oldUrl := c.baseUrl
	c.baseUrl = sessionserverUrl // i dont like this but i cant think of any other way atm :(
	endpoint := "session/minecraft/profile/" + uuid

	response, err := get[UuidToProfileResponse](c, ctx, endpoint, http.Header{})
	c.baseUrl = oldUrl
	return response, err
}
