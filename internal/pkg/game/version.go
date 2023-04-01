package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

const (
	versionManifestUrl  = "https://launchermeta.mojang.com/mc/game/version_manifest.json"
	versionManifestFile = "versions.json"
)

type (
	VersionManifest struct {
		LastUpdated time.Time `json:"lastUpdated"`
		Latest      struct {
			Release  string `json:"release"`
			Snapshot string `json:"snapshot"`
		}
		Versions []*VersionInfo
	}
	VersionInfo struct {
		Id          string    `json:"id"`
		ReleaseTime time.Time `json:"releaseTime"`
		Time        time.Time `json:"time"`
		Type        string    `json:"type"`
		Url         string    `json:"url"`
	}
)

// Version manager

type VersionManager struct {
	Manifest *VersionManifest
}

func NewVersionManager(dataDir string) (*VersionManager, error) {
	versionsFile := path.Join(dataDir, versionManifestFile)

	var manifest VersionManifest
	if _, err := os.Stat(versionsFile); errors.Is(err, fs.ErrNotExist) {
		//todo downloading the manifest should show a loading indicator in the cli
		println("downloading manifest")

		// Download the manifest
		manifest, err = downloadVersionManifest(versionsFile)
		if err != nil {
			return nil, err
		}
	} else {
		// Load the manifest
		f, err := os.Open(versionsFile)
		if err != nil {
			return nil, fmt.Errorf("failed to open version manifest file: %w", err)
		}
		defer f.Close()

		if err := json.NewDecoder(f).Decode(&manifest); err != nil {
			return nil, fmt.Errorf("failed to read version manifest: %w", err)
		}
	}

	return &VersionManager{
		Manifest: &manifest,
	}, nil
}

func (m *VersionManager) FindVersionByName(name string) *VersionInfo {
	name = strings.ToLower(name)
	for _, version := range m.Manifest.Versions {
		if version.Id == name {
			return version
		}
	}
	return nil
}

func downloadVersionManifest(path string) (manifest VersionManifest, err error) {
	res, err := http.Get(versionManifestUrl)
	if err != nil {
		return manifest, fmt.Errorf("failed to get version manifest: %w", err)
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&manifest); err != nil {
		return manifest, fmt.Errorf("failed to Download version manifest: %w", err)
	}
	manifest.LastUpdated = time.Now()

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return manifest, fmt.Errorf("failed to open version manifest file: %w", err)
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(&manifest); err != nil {
		return manifest, fmt.Errorf("failed to write version manifest file: %w", err)
	}

	return manifest, nil
}
