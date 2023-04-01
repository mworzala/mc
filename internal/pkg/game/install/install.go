package install

import (
	"errors"
	"fmt"
	"path"
	"strings"
	"sync"

	gameModel "github.com/mworzala/mc-cli/internal/pkg/game/model"
	"github.com/mworzala/mc-cli/internal/pkg/game/rule"
	"github.com/mworzala/mc-cli/internal/pkg/util"
)

var (
	MinLauncherVersion            = 21
	MaxLauncherVersion            = 21
	ErrUnsupportedLauncherVersion = errors.New("unsupported launcher version")
)

type Installer struct {
	configDir      string
	getVersionFunc func(string) (*gameModel.VersionInfo, error)

	// Common directories
	versionsDir  string
	librariesDir string
	assetsDir    string

	rules *rule.Evaluator
}

func NewInstaller(configDir string, getVersionFunc func(string) (*gameModel.VersionInfo, error)) *Installer {
	return &Installer{
		configDir:      configDir,
		getVersionFunc: getVersionFunc,

		versionsDir:  path.Join(configDir, "versions"),
		librariesDir: path.Join(configDir, "libraries"),
		assetsDir:    path.Join(configDir, "assets"),

		rules: rule.NewEvaluator(),
	}
}

func (i *Installer) Install(v *gameModel.VersionInfo) error {
	versionDir := path.Join(i.versionsDir, v.Id)

	// Download the version spec (or read it if it exists)
	var spec gameModel.VersionSpec
	versionSpecPath := path.Join(versionDir, fmt.Sprintf("%s.json", v.Id))
	if err := util.ReadOrDownload(v.Id, versionSpecPath, util.FileDownload{Url: v.Url}, &spec); err != nil {
		return fmt.Errorf("failed to read version spec: %w", err)
	}

	// If there is an inherited version, install that
	if spec.InheritsFrom != "" {
		inherited, err := i.getVersionFunc(spec.InheritsFrom)
		if err != nil {
			return fmt.Errorf("inherited version not found: %s", spec.InheritsFrom)
		}

		if err := i.Install(inherited); err != nil {
			return fmt.Errorf("error installing inherited version %s: %w", spec.InheritsFrom, err)
		}
	}

	return i.InstallFromSpec(&spec)
}

func (i *Installer) InstallFromSpec(spec *gameModel.VersionSpec) error {
	// We assume support if the version is zero - fabric does not provide a version
	if spec.MinimumLauncherVersion != 0 &&
		(spec.MinimumLauncherVersion < MinLauncherVersion ||
			spec.MinimumLauncherVersion > MaxLauncherVersion) {
		return fmt.Errorf("%w: %d", ErrUnsupportedLauncherVersion, spec.MinimumLauncherVersion)
	}

	// Download client archive
	if spec.Downloads != nil && spec.Downloads.Client != nil {
		clientPath := path.Join(i.versionsDir, spec.Id, fmt.Sprintf("%s.jar", spec.Id))
		if err := util.Download("", clientPath, *spec.Downloads.Client); err != nil {
			return fmt.Errorf("failed to download client: %w", err)
		}
	}

	// Libraries
	if err := i.downloadLibraries(spec.Libraries); err != nil {
		return err
	}

	// Asset index
	if index := spec.AssetIndex; index != nil {
		var assetIndex gameModel.AssetIndex
		assetIndexPath := path.Join(i.assetsDir, "indexes", fmt.Sprintf("%s.json", index.Id))
		if err := util.ReadOrDownload("", assetIndexPath, index.FileDownload, &assetIndex); err != nil {
			return fmt.Errorf("failed to download asset index: %w", err)
		}

		// Asset objects
		if err := i.downloadAssetObjects(index.TotalSize, &assetIndex); err != nil {
			return err
		}
	}

	// Log config
	if logging := spec.Logging; logging != nil {
		logConfigPath := path.Join(i.assetsDir, "log_configs", logging.Client.File.Id)
		if err := util.Download("", logConfigPath, logging.Client.File.FileDownload); err != nil {
			return fmt.Errorf("failed to download log config: %w", err)
		}
	}

	return nil
}

func (i *Installer) downloadLibraries(libraries []*gameModel.Library) error {
	for _, library := range libraries {
		if i.rules.Eval(library.Rules) == rule.Deny {
			continue
		}

		if library.Downloads != nil { // Vanilla-type library
			artifact := library.Downloads.Artifact
			libraryPath := path.Join(i.librariesDir, artifact.Path)
			if err := util.Download("", libraryPath, artifact.FileDownload); err != nil {
				return fmt.Errorf("failed to download library %s: %w", library.Name, err)
			}
		} else if library.Url != "" { // Direct maven library
			parts := strings.Split(library.Name, ":")
			groupId := parts[0]
			artifactName := parts[1]
			version := parts[2]

			artifactPath := fmt.Sprintf("%s/%s/%s/%s-%s.jar", strings.ReplaceAll(groupId, ".", "/"), artifactName, version, artifactName, version)
			artifactUrl := fmt.Sprintf("%s/%s", strings.TrimSuffix(library.Url, "/"), artifactPath)

			if err := util.Download("", path.Join(i.librariesDir, artifactPath), util.FileDownload{Url: artifactUrl}); err != nil {
				return fmt.Errorf("failed to download library %s: %w", library.Name, err)
			}
		}
	}
	return nil
}

func (i *Installer) downloadAssetObjects(totalSize int64, index *gameModel.AssetIndex) error {
	openConns := make(chan struct{}, 150)
	for i := 0; i < 150; i++ {
		openConns <- struct{}{}
	}

	objectsPath := path.Join(i.assetsDir, "objects")

	wg := sync.WaitGroup{}
	wg.Add(len(index.Objects))
	for _, obj := range index.Objects {
		go func(obj *gameModel.AssetObject) {
			defer wg.Done()

			// Read from connection pool and then at the end write back to it
			<-openConns
			defer func() {
				openConns <- struct{}{}
			}()

			objPath := path.Join(objectsPath, obj.Hash[:2], obj.Hash)
			objUrl := fmt.Sprintf("%s/%s/%s", gameModel.MojangObjectBaseUrl, obj.Hash[:2], obj.Hash)

			dl := util.FileDownload{Sha1: obj.Hash, Size: obj.Size, Url: objUrl}
			if err := util.Download("", objPath, dl); err != nil {
				panic(err) //todo handle this case better
			}
		}(obj)
	}

	wg.Wait()
	return nil

}
