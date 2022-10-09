package install

import (
	"errors"
	"fmt"
	"github.com/mworzala/mc-cli/internal/pkg/model"
	"github.com/mworzala/mc-cli/internal/pkg/util"
	"os"
	"path"
	"sync"
)

func Vanilla(ctx util.Context, version *model.ManifestVersion) error {
	id := version.Id
	versionDir := path.Join(ctx.RootDir, "versions", id)

	// Download the version manifest
	//todo just remove the progress bar for this download
	var v model.Version
	versionDescPath := path.Join(versionDir, fmt.Sprintf("%s.json", id))
	err := ctx.ReadOrDownload("id", versionDescPath, model.Download{Url: version.Url}, &v)
	if err != nil {
		return err
	}

	//todo validate the minimumManifestVersion

	// Download the client jar
	clientJarPath := path.Join(versionDir, fmt.Sprintf("%s.jar", id))
	err = ctx.Download(fmt.Sprintf("%s.jar", id), clientJarPath, v.Downloads.Client)
	if err != nil {
		return err
	}

	// Download libraries
	librariesPath := path.Join(ctx.RootDir, "libraries")
	err = downloadVanillaLibraries(ctx, librariesPath, v.Libraries)

	// Download asset metadata
	assetPath := path.Join(ctx.RootDir, "assets")

	// Log config
	logConfig := v.Logging.Client
	logConfigPath := path.Join(assetPath, "log_configs", logConfig.File.Id)
	err = ctx.Download(logConfig.File.Id, logConfigPath, logConfig.File.Download)

	// Asset Index
	assets := v.AssetIndex
	var assetIndex model.AssetIndex
	assetIndexPath := path.Join(assetPath, "indexes", fmt.Sprintf("%s.json", assets.Id))
	err = ctx.ReadOrDownload(fmt.Sprintf("%s.json", assets.Id), assetIndexPath, assets.Download, &assetIndex)

	// Objects
	objectsPath := path.Join(assetPath, "objects")
	err = downloadVanillaObjects(ctx, objectsPath, assets.TotalSize, assetIndex)
	if err != nil {
		return err
	}

	return nil
}

func downloadVanillaLibraries(ctx util.Context, librariesPath string, libraries []*model.Library) error {
	for _, library := range libraries {
		// Check rules
		if len(library.Rules) > 0 {
			allowed := false
			for _, rule := range library.Rules {
				allowed = EvalRule(rule)
			}
			if !allowed {
				continue
			}
		}

		// Download artifact
		artifact := library.Downloads.Artifact
		libPath := path.Join(librariesPath, artifact.Path)
		err := ctx.Download(library.Name, libPath, artifact.Download)
		if err != nil {
			return fmt.Errorf("failed to download library %s: %w", library.Name, err)
		}
	}
	return nil
}

func EvalRule(rule map[string]interface{}) bool {
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

func downloadVanillaObjects(ctx util.Context, objectsPath string, totalSize int64, assetIndex model.AssetIndex) error {
	// Create single progress bar for all downloads
	fmt.Println("\nassets")
	bar := ctx.CreateProgressBar(totalSize)

	// Download all
	wg := sync.WaitGroup{}
	wg.Add(len(assetIndex.Objects))
	for _, obj := range assetIndex.Objects {
		obj := obj
		go func() {
			objPath := path.Join(objectsPath, obj.Hash[:2], obj.Hash)
			objUrl := fmt.Sprintf("https://resources.download.minecraft.net/%s/%s", obj.Hash[:2], obj.Hash)

			if _, err := os.Stat(objPath); errors.Is(err, os.ErrNotExist) {
				dl := model.Download{Sha1: obj.Hash, Size: obj.Size, Url: objUrl}
				err = ctx.DownloadFile(objPath, dl, bar)
				if err != nil {
					panic(err) //todo handle this case better
				}
			} else {
				_ = bar.Add64(obj.Size)
			}

			wg.Done()
		}()
	}

	wg.Wait()
	return nil
}
