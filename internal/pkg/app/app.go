package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mworzala/mc-cli/internal/pkg/auth"
	"github.com/mworzala/mc-cli/internal/pkg/install"
	"github.com/mworzala/mc-cli/internal/pkg/model"
	"github.com/mworzala/mc-cli/internal/pkg/util"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

type App struct {
	rootDir      string
	authProvider auth.Provider
}

func NewApp(rootDir string) *App {
	return &App{
		rootDir:      rootDir,
		authProvider: auth.NewProvider(rootDir),
	}
}

func (a *App) ReadManifest() (*model.Manifest, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ctx := a.CreateContextFrom(timeoutCtx)

	var manifest model.Manifest
	err := ctx.ReadOrDownload("", path.Join(a.rootDir, "manifest.json"), model.Download{
		Url: "https://launchermeta.mojang.com/mc/game/version_manifest.json",
	}, &manifest)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest: %w", err)
	}

	return &manifest, nil
}

func (a *App) RunLatest() error {
	manifest, err := a.ReadManifest()
	if err != nil {
		return err
	}

	// Find latest version
	var latest *model.ManifestVersion
	for _, v := range manifest.Versions {
		if v.Id == manifest.Latest.Release {
			latest = v
		}
	}

	// Get latest credentials
	uuid := "aceb326fda1545bcbf2f11940c21780c"
	credentials, err := a.authProvider.GetCredentials(uuid)
	if err != nil {
		if err == auth.ErrNotFound {
			fmt.Printf("no credentials found for %s\n", uuid)
			return nil
		}
		//todo other errors such as unable to refresh token requesting that they sign in again

		return fmt.Errorf("failed to get credentials: %w", err)
	}
	fmt.Printf("🔒 Authenticated as %s\n", credentials.PlayerName)

	// Download if missing
	ok, err := a.IsInstalled(latest.Id)
	if err != nil {
		return fmt.Errorf("unexpected error: %w", err)
	}
	if !ok {
		fmt.Printf("📁 Installing %s\n", latest.Id)
		ctx := a.CreateContext()
		err := install.Vanilla(ctx, latest)
		if err != nil {
			return fmt.Errorf("failed to install %s version: %w", latest.Id, err)
		}
	}

	// Launch
	fmt.Printf("🚀 Launching %s\n", latest.Id)

	//todo cleanup
	f, err := os.OpenFile(path.Join(a.rootDir, "versions", latest.Id, fmt.Sprintf("%s.json", latest.Id)), os.O_RDONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open version json: %w", err)
	}
	defer f.Close()
	var v model.Version
	err = json.NewDecoder(f).Decode(&v)
	if err != nil {
		return fmt.Errorf("failed to decode version json: %w", err)
	}

	var args []string
	classpath := strings.Builder{}
	librariesPath := path.Join(a.rootDir, "libraries")
	for _, lib := range v.Libraries {
		libPath := path.Join(librariesPath, lib.Downloads.Artifact.Path)
		classpath.WriteString(libPath)
		classpath.WriteString(":")
	}
	classpath.WriteString(path.Join(a.rootDir, "versions", v.Id, v.Id+".jar"))

	vars := map[string]string{
		// jvm
		"natives_directory": ".",
		"launcher_name":     "mc-cli",
		"launcher_version":  "0.0.1",
		"classpath":         classpath.String(),
		// game
		"version_name":      latest.Id,
		"game_directory":    path.Join(a.rootDir, "tmp-instance"),
		"assets_root":       path.Join(a.rootDir, "assets"),
		"assets_index_name": v.Assets,
		"auth_player_name":  credentials.PlayerName,
		"auth_uuid":         credentials.UUID,
		"auth_access_token": credentials.AccessToken,
		"clientid":          "MTMwQUU2ODYwQUE1NDUwNkIyNUZCMzZBNjFCNjc3M0Q=",
		"auth_xuid":         credentials.UserHash,
		"user_type":         string(credentials.UserType),
		"version_type":      latest.Type,
	}
	_ = vars

	replaceVars := func(s string) string {
		for k, v := range vars {
			s = strings.ReplaceAll(s, fmt.Sprintf("${%s}", k), v)
		}
		return s
	}

	args = append(args, "-XstartOnFirstThread")

	for _, arg := range v.Arguments.JVM {
		if s, ok := arg.(string); ok {
			args = append(args, replaceVars(s))
		} else if m, ok := arg.(map[string]interface{}); ok {
			_ = m
			//value := m["value"]
			//if s, ok := value.(string); ok {
			//	args = append(args, replaceVars(s))
			//} else if a, ok := value.([]interface{}); ok {
			//	for _, v := range a {
			//		if s, ok := v.(string); ok {
			//			args = append(args, replaceVars(s))
			//		}
			//	}
			//} else {
			//	panic(fmt.Sprintf("unknown type: %T", value))
			//}
		} else {
			panic("unknown arg type")
		}
	}

	args = append(args, v.MainClass)

	for _, arg := range v.Arguments.Game {
		if s, ok := arg.(string); ok {
			args = append(args, replaceVars(s))
		} else if m, ok := arg.(map[string]interface{}); ok {
			_ = m
			//value := m["value"]
			//if s, ok := value.(string); ok {
			//	args = append(args, replaceVars(s))
			//} else if a, ok := value.([]interface{}); ok {
			//	for _, v := range a {
			//		if s, ok := v.(string); ok {
			//			args = append(args, replaceVars(s))
			//		}
			//	}
			//} else {
			//	panic(fmt.Sprintf("unknown type: %T", value))
			//}
		} else {
			panic("unknown arg type")
		}
	}

	javaBin := "/Users/matt/Library/Java/JavaVirtualMachines/liberica-17.0.1/bin/java"
	cmd := exec.Command(javaBin, args...)
	cmd.Dir = path.Join(a.rootDir, "tmp-instance")

	cmd.Stdout = os.Stdout

	if err != nil {
		panic(err)
	}
	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	err = cmd.Wait()
	if err != nil {
		panic(err)
	}

	return nil
}

func (a *App) LoginMicrosoft() error {
	user, err := a.authProvider.LoginMSA()
	if err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}

	fmt.Printf("🎉 Signed in as %s\n", user)
	return nil
}

// Helpers

func (a *App) CreateContext() util.Context {
	return util.Context{
		RootDir: a.rootDir,
		Ctx:     context.Background(),
	}
}

func (a *App) CreateContextFrom(ctx context.Context) util.Context {
	return util.Context{
		RootDir: a.rootDir,
		Ctx:     ctx,
	}
}

func (a *App) IsInstalled(version string) (bool, error) {
	versionDir := path.Join(a.rootDir, "versions", version)
	// todo this check is not really valid, there are some error cases.
	if _, err := os.Stat(versionDir); !errors.Is(err, os.ErrNotExist) {
		return true, nil
	}
	return false, nil
}

//func downloadFileIfNotExists(id, file string, dl model.Download) error {
//	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
//		resp, err := http.Get(dl.Url)
//		if err != nil {
//			return err
//		}
//		defer resp.Body.Close()
//
//		//todo breaks for asset index for some reason
//		//if resp.ContentLength != dl.Size {
//		//	return fmt.Errorf("download size mismatch for %s: %d != %d", id, resp.ContentLength, dl.Size)
//		//}
//
//		f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0644)
//		if err != nil {
//			return err
//		}
//		defer f.Close()
//
//		bar := progressbar.DefaultBytes(resp.ContentLength, id)
//
//		h := sha1.New()
//		_, err = io.Copy(io.MultiWriter(f, h, bar), resp.Body)
//		if err != nil {
//			return err
//		}
//
//		hash := fmt.Sprintf("%x", h.Sum(nil))
//		if hash != dl.Sha1 {
//			return fmt.Errorf("download hash mismatch for %s: %s != %s", id, hash, dl.Sha1)
//		}
//	}
//
//	return nil
//}
//
//func getVersion(manifest *model.Manifest, version string) *model.ManifestVersion {
//	for _, v := range manifest.Versions {
//		if v.Id == version {
//			return v
//		}
//	}
//	return nil
//}
//
//func writeManifest(manifest *model.Manifest) error {
//	file, err := os.OpenFile(path.Join(basePath, "manifest.json"), os.O_CREATE|os.O_WRONLY, 0644)
//	if err != nil {
//		return err
//	}
//	defer file.Close()
//
//	err = json.NewEncoder(file).Encode(manifest)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func writeJson(file string, data interface{}) error {
//	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0644)
//	if err != nil {
//		return err
//	}
//	defer f.Close()
//
//	err = json.NewEncoder(f).Encode(data)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
