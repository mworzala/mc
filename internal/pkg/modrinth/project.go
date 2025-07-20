package modrinth

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

var (
	ErrHashMismatch = errors.New("download hash was mismatched")
)

type DependencyType string

const (
	DependencyRequired DependencyType = "required"
)

type ProjectVersion struct {
	Id           string              `json:"id"`
	ProjectId    string              `json:"project_id"`
	Files        []VersionFile       `json:"files"`
	Dependencies []VersionDependency `json:"dependencies"`
}

type VersionFile struct {
	Hashes struct {
		Sha1   string `json:"sha1"`
		Sha512 string `json:"sha512"`
	} `json:"hashes"`
	Url      string  `json:"url"`
	Filename string  `json:"filename"`
	Primary  bool    `json:"primary"`
	Size     int     `json:"size"`
	FileType *string `json:"file_type"`
}

type VersionDependency struct {
	VersionId      *string        `json:"version_id"`
	ProjectId      *string        `json:"project_id"`
	DependencyType DependencyType `json:"dependency_type"`
}

func (c *Client) GetProjectVersions(ctx context.Context, modId string, gameVersions []string, loaders []string) (*[]ProjectVersion, error) {
	jsonGameVersions, err := json.Marshal(gameVersions)
	if err != nil {
		return nil, err
	}

	jsonLoaders, err := json.Marshal(loaders)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Add("game_versions", string(jsonGameVersions))
	params.Add("loaders", string(jsonLoaders))

	return get[[]ProjectVersion](c, ctx, fmt.Sprintf("/project/%s/version", modId), params)
}

func (c *Client) GetProjectVersion(ctx context.Context, versionId string) (*ProjectVersion, error) {
	return get[ProjectVersion](c, ctx, fmt.Sprintf("/version/%s", versionId), url.Values{})
}

func (c *Client) DownloadFile(ctx context.Context, file VersionFile, tempFile *os.File) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, file.Url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", c.userAgent)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusBadRequest {
		var errorRes badRequestError
		if err := json.NewDecoder(res.Body).Decode(&errorRes); err != nil {
			return fmt.Errorf("failed to decode response body: %w", err)
		}

		return fmt.Errorf("400 %s: %s", errorRes.Error, errorRes.Description)
	} else if res.StatusCode == http.StatusInternalServerError {
		return errors.New("500 internal server error")
	} else if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response from server: %d", res.StatusCode)
	}

	hasher := sha512.New()
	writer := io.MultiWriter(tempFile, hasher)

	if _, err := io.Copy(writer, res.Body); err != nil {
		return err
	}

	hash := hex.EncodeToString(hasher.Sum(nil))
	if hash != file.Hashes.Sha512 {
		return ErrHashMismatch
	}

	return nil
}
