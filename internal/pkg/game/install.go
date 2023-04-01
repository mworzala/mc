package game

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

type VersionSpec struct {
	Arguments struct {
		Game []interface{} `json:"game"`
		JVM  []interface{} `json:"jvm"`
	} `json:"arguments"`
	AssetIndex struct {
		Id        string `json:"id"`
		TotalSize int64  `json:"totalSize"`
		Download
	} `json:"assetIndex"`
	Assets string `json:"assets"`
	// ComplianceLevel of 1 indicates that it supports "new security features"
	// Its used for that warning in the launcher.
	ComplianceLevel int `json:"complianceLevel"`
	Downloads       struct {
		Client         Download `json:"client"`
		ClientMappings Download `json:"client_mappings"`
		Server         Download `json:"server"`
		ServerMappings Download `json:"server_mappings"`
	} `json:"downloads"`
	Id          string `json:"id"`
	JavaVersion struct {
		Component    string `json:"component"`
		MajorVersion int    `json:"majorVersion"`
	} `json:"javaVersion"`
	Libraries []*Library `json:"libraries"`
	Logging   struct {
		Client struct {
			Argument string `json:"argument"`
			File     struct {
				Id string `json:"id"`
				Download
			} `json:"file"`
			Type string `json:"type"`
		} `json:"client"`
	} `json:"logging"` //todo
	MainClass              string    `json:"mainClass"`
	MinimumLauncherVersion int       `json:"minimumLauncherVersion"`
	ReleaseTime            time.Time `json:"releaseTime"`
	Time                   time.Time `json:"time"`
	Type                   string    `json:"type"`
}

type Library struct {
	Name      string `json:"name"`
	Downloads struct {
		Artifact struct {
			Path string `json:"path"`
			Download
		} `json:"artifact"`
	} `json:"downloads"`
	Rules []map[string]interface{} `json:"rules"` //todo
}

type AssetIndex struct {
	Objects map[string]*AssetObject `json:"objects"`
}

type AssetObject struct {
	Hash string `json:"hash"`
	Size int64  `json:"size"`
}

type Download struct {
	Sha1 string `json:"sha1"`
	Size int64  `json:"size"`
	Url  string `json:"url"`
}

func InstallVersion(dataDir string, v *VersionInfo) error {
	versionDir := path.Join(dataDir, "versions", v.Id)

	// Download the version manifest
	var spec VersionSpec
	versionSpecPath := path.Join(versionDir, fmt.Sprintf("%s.json", v.Id))
	if err := readOrDownload("id", versionSpecPath, Download{Url: v.Url}, &spec); err != nil {
		return err
	}

	//todo validate the minimumManifestVersion

	// Download the client jar
	clientJarPath := path.Join(versionDir, fmt.Sprintf("%s.jar", v.Id))
	if err := download(path.Base(clientJarPath), clientJarPath, spec.Downloads.Client); err != nil {
		return err
	}

	// Download libraries
	librariesPath := path.Join(dataDir, "libraries")
	if err := downloadLibraries(librariesPath, spec.Libraries); err != nil {
		return err
	}

	// Download asset index
	assetPath := path.Join(dataDir, "assets")

	var assetIndex AssetIndex
	assetIndexPath := path.Join(assetPath, "indexes", fmt.Sprintf("%s.json", spec.AssetIndex.Id))
	if err := readOrDownload(path.Base(assetIndexPath), assetIndexPath, spec.AssetIndex.Download, &assetIndex); err != nil {
		return err
	}

	// Download objects (from asset index)
	objectsPath := path.Join(assetPath, "objects")
	if err := downloadObjects(objectsPath, spec.AssetIndex.TotalSize, &assetIndex); err != nil {
		return err
	}

	// Download log config
	logConfigSpec := spec.Logging.Client
	logConfigPath := path.Join(assetPath, "log_configs", logConfigSpec.File.Id)
	if err := download(logConfigSpec.File.Id, logConfigPath, logConfigSpec.File.Download); err != nil {
		return err
	}

	return nil
}

func downloadLibraries(librariesPath string, libraries []*Library) error {
	for _, library := range libraries {
		// Check rules
		if len(library.Rules) > 0 {
			allowed := false
			for _, rule := range library.Rules {
				allowed = evalRule(rule)
			}
			if !allowed {
				continue
			}
		}

		// Download artifact
		artifact := library.Downloads.Artifact
		libPath := path.Join(librariesPath, artifact.Path)
		if err := download(library.Name, libPath, artifact.Download); err != nil {
			return fmt.Errorf("failed to download library %s: %w", library.Name, err)
		}
	}

	return nil
}

func evalRule(rule map[string]interface{}) bool {
	action, ok := rule["action"].(string)
	if !ok {
		panic("invalid rule") //todo
	}
	value := action == "allow"

	if os_, ok := rule["os"].(map[string]interface{}); ok {
		if name, ok := os_["name"].(string); ok && name != "osx" {
			return !value
		}
		if arch, ok := os_["arch"].(string); ok && arch != "arm64" {
			return !value
		}
	}
	return value
}

func downloadObjects(objectsPath string, totalSize int64, assetIndex *AssetIndex) error {
	openConns := make(chan struct{}, 50)
	for i := 0; i < 50; i++ {
		openConns <- struct{}{}
	}

	wg := sync.WaitGroup{}
	wg.Add(len(assetIndex.Objects))
	for _, obj := range assetIndex.Objects {
		go func(obj *AssetObject) {
			defer wg.Done()

			// Read from connection pool and then at the end write back to it
			<-openConns
			defer func() {
				openConns <- struct{}{}
			}()

			objPath := path.Join(objectsPath, obj.Hash[:2], obj.Hash)
			objUrl := fmt.Sprintf("https://resources.download.minecraft.net/%s/%s", obj.Hash[:2], obj.Hash)

			if _, err := os.Stat(objPath); errors.Is(err, fs.ErrNotExist) {
				dl := Download{Sha1: obj.Hash, Size: obj.Size, Url: objUrl}
				if err := downloadFile(objPath, dl); err != nil {
					panic(err) //todo handle this case better
				}
			}
		}(obj)
	}

	wg.Wait()
	return nil
}

func readOrDownload(id, file string, download Download, ptr interface{}) error {
	if _, err := os.Stat(file); err == nil {
		return readFile(file, ptr)
	} else if errors.Is(err, fs.ErrNotExist) {
		data := new(bytes.Buffer)

		println("Download", id) //todo remove me/add some callback for progress
		if err := downloadFile(file, download, data); err != nil {
			return err
		}

		return json.NewDecoder(data).Decode(ptr)
	} else {
		return err
	}
}

func download(id, file string, dl Download) error {
	if _, err := os.Stat(file); errors.Is(err, fs.ErrNotExist) {
		println("Download", id) //todo remove me/add some callback for progress
		return downloadFile(file, dl)
	}
	return nil
}

func downloadFile(file string, download Download, listeners ...io.Writer) error {
	// Create parent directory
	if err := os.MkdirAll(path.Dir(file), 0755); err != nil {
		return err
	}

	// Open request
	res, err := http.Get(download.Url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Create file
	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	// Copy data to file, hash, and listeners
	hash := sha1.New()
	writers := append(listeners, f)
	if download.Sha1 != "" {
		writers = append(writers, hash)
	}
	if _, err := io.Copy(io.MultiWriter(writers...), res.Body); err != nil {
		return err
	}

	// Validate hash if present
	if download.Sha1 != "" {
		h := fmt.Sprintf("%x", hash.Sum(nil))
		if h != download.Sha1 {
			// Attempt to delete the file
			_ = os.Remove(file)

			return fmt.Errorf("Download hash mismatch: %s != %s", h, download.Sha1)
		}
	}

	return nil
}

func readFile(file string, ptr interface{}) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(ptr)
}
