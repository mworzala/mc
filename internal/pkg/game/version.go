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

	gameModel "github.com/mworzala/mc/internal/pkg/game/model"
)

const (
	// Vanilla
	versionManifestUrl             = "https://launchermeta.mojang.com/mc/game/version_manifest_v2.json"
	experimentalVersionManifestUrl = "https://maven.fabricmc.net/net/minecraft/experimental_versions.json"

	// Fabric
	fabricVersionManifestUrl = "https://meta.fabricmc.net/v2/versions/game"
	fabricLoaderManifestUrl  = "https://meta.fabricmc.net/v2/versions/loader"
	fabricVersionSpecBaseUrl = "https://meta.fabricmc.net/v2/versions/loader"

	versionManifestV2File = "versions_v2.json"
)

var (
	ErrUnknownVersion       = errors.New("unknown version")
	ErrUnknownFabricVersion = errors.New("unknown fabric version")
	ErrUnknownFabricLoader  = errors.New("unknown fabric loader")
	TriedToUpdate           = false
)

type (
	VersionManifestV2 struct {
		LastUpdated time.Time
		Vanilla     struct {
			Release  string
			Snapshot string
			// Mapping of Minecraft version to version json url
			Versions map[string]*gameModel.VersionInfo
		}
		Fabric struct {
			// Versions matrix is a mapping of Minecraft version to fabric support
			// The entries in here are partial, and should not be used as is
			Versions      map[string]*gameModel.VersionInfo
			DefaultLoader string
			Loaders       map[string]bool
		}
	}

	// Response from versionManifestUrl and experimentalVersionManifestUrl
	mojangVersionManifestV2 struct {
		Latest *struct {
			Release  string `json:"release"`
			Snapshot string `json:"snapshot"`
		}
		Versions []struct {
			Id          string    `json:"id"`
			ReleaseTime time.Time `json:"releaseTime"`
			Time        time.Time `json:"time"`
			Type        string    `json:"type"`
			Url         string    `json:"url"`
		} `json:"versions"`
	}
	// Response from fabricVersionManifestUrl
	fabricVersionManifestV2 []struct {
		Version string `json:"version"`
		Stable  bool   `json:"stable"`
	}
	// Response from fabricLoaderManifestUrl
	fabricLoaderManifestV2 []struct {
		Separator string `json:"separator"`
		Build     int    `json:"build"`
		Maven     string `json:"maven"`
		Version   string `json:"version"`
		Stable    bool   `json:"stable"`
	}
)

// Version manager

type VersionManager struct {
	cacheFile  string
	manifestV2 *VersionManifestV2
}

func NewVersionManager(dataDir string) (*VersionManager, error) {
	cacheFile := path.Join(dataDir, versionManifestV2File)
	m := &VersionManager{cacheFile: cacheFile}

	if _, err := os.Stat(cacheFile); errors.Is(err, fs.ErrNotExist) {
		if err := m.updateManifest(); err != nil {
			return nil, err
		}
	} else {
		f, err := os.Open(cacheFile)
		if err != nil {
			return nil, fmt.Errorf("failed to open version manifest file: %w", err)
		}
		defer f.Close()

		var manifest VersionManifestV2
		if err := json.NewDecoder(f).Decode(&manifest); err != nil {
			return nil, fmt.Errorf("failed to read version manifest: %w", err)
		}
		m.manifestV2 = &manifest
	}

	//todo update if old

	return m, nil
}

func (m *VersionManager) FindVanilla(name string) (*gameModel.VersionInfo, error) {
	v, ok := m.manifestV2.Vanilla.Versions[strings.ToLower(name)]
	if !ok {
		if TriedToUpdate {
			return nil, ErrUnknownVersion
		}
		TriedToUpdate = true
		fmt.Println("Couldn't find version: ", name, ", refreshing manfiest")
		m.updateManifest()
		return m.FindVanilla(name)
	}
	TriedToUpdate = false
	return v, nil
}

func (m *VersionManager) FindFabric(name, loader string) (*gameModel.VersionInfo, error) {
	if !m.FabricLoaderExists(loader) {
		return nil, ErrUnknownFabricLoader
	}

	partial, ok := m.manifestV2.Fabric.Versions[strings.ToLower(name)]
	if !ok {
		if TriedToUpdate {
			return nil, ErrUnknownFabricVersion
		}
		TriedToUpdate = true
		fmt.Println("Couldn't find fabric for version: ", name, ", refreshing manfiest")
		m.updateManifest()
		return m.FindFabric(name, loader)
	}
	TriedToUpdate = false

	return &gameModel.VersionInfo{
		Id:     fmt.Sprintf(partial.Id, loader),
		Stable: partial.Stable,
		Url:    fmt.Sprintf(partial.Url, loader),
	}, nil
}

func (m *VersionManager) DefaultFabricLoader() string {
	return m.manifestV2.Fabric.DefaultLoader
}

func (m *VersionManager) FabricLoaderExists(name string) bool {
	_, ok := m.manifestV2.Fabric.Loaders[strings.ToLower(name)]
	return ok
}

func (m *VersionManager) updateManifest() error {
	var result VersionManifestV2
	result.LastUpdated = time.Now()
	result.Vanilla.Versions = make(map[string]*gameModel.VersionInfo)
	result.Fabric.Versions = make(map[string]*gameModel.VersionInfo)
	result.Fabric.Loaders = make(map[string]bool)

	updateMojangManifest := func(url string) error {
		res, err := http.Get(url)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		var manifest mojangVersionManifestV2
		if err := json.NewDecoder(res.Body).Decode(&manifest); err != nil {
			return err
		}

		if manifest.Latest != nil {
			result.Vanilla.Release = manifest.Latest.Release
			result.Vanilla.Snapshot = manifest.Latest.Snapshot
		}
		for _, v := range manifest.Versions {
			result.Vanilla.Versions[v.Id] = &gameModel.VersionInfo{
				Id:     v.Id,
				Stable: v.Type == "release",
				Url:    v.Url,
			}
		}

		return nil
	}

	// Pull vanilla and experimental manifests
	if err := updateMojangManifest(versionManifestUrl); err != nil {
		return fmt.Errorf("failed to update mojang manifest: %w", err)
	}
	if err := updateMojangManifest(experimentalVersionManifestUrl); err != nil {
		return fmt.Errorf("failed to update experimental manifest: %w", err)
	}

	// Pull fabric loader manifest
	{
		res, err := http.Get(fabricLoaderManifestUrl)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		var manifest fabricLoaderManifestV2
		if err := json.NewDecoder(res.Body).Decode(&manifest); err != nil {
			return err
		}

		for _, v := range manifest {
			if v.Stable && result.Fabric.DefaultLoader == "" {
				result.Fabric.DefaultLoader = v.Version
			}
			result.Fabric.Loaders[v.Version] = v.Stable
		}
	}

	// Pull fabric versions
	{
		res, err := http.Get(fabricVersionManifestUrl)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		var manifest fabricVersionManifestV2
		if err := json.NewDecoder(res.Body).Decode(&manifest); err != nil {
			return err
		}

		for _, v := range manifest {
			result.Fabric.Versions[v.Version] = &gameModel.VersionInfo{
				Id:     fmt.Sprintf("fabric-loader-%%s-%s", v.Version),
				Stable: v.Stable,
				Url:    fmt.Sprintf("%s/%s/%%s/profile/json", fabricVersionSpecBaseUrl, v.Version),
			}
		}
	}

	f, err := os.OpenFile(m.cacheFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(&result); err != nil {
		return err
	}

	m.manifestV2 = &result
	return nil
}
