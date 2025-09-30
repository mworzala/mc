package profile

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mworzala/mc/internal/pkg/modrinth"
)

type Provider string

const (
	mcCliDir     = "/.mc-cli/"
	modsFileName = "mods.json"

	Modrinth Provider = "modrinth"
	//CurseForge Provider = "curseforge" // Requires an API key from what I can tell, ignoring for now
)

var (
	ErrModAlreadyInstalled = errors.New("mod already installed")
)

type Mod struct {
	Id        string   `json:"id"`
	VersionId string   `json:"version_id"`
	Enabled   bool     `json:"enabled"`
	Provider  Provider `json:"provider"`
	JarFile   string   `json:"jar_file"`
}

type UnsavedMod struct {
	Id        string               `json:"id"`
	VersionId string               `json:"version_id"`
	Provider  Provider             `json:"provider"`
	File      modrinth.VersionFile `json:"file"`
}

func (p *Profile) Mods() ([]Mod, error) {
	if p.mods != nil {
		return p.mods, nil
	}

	mccliDir := filepath.Join(p.Directory, mcCliDir)
	name := filepath.Join(mccliDir, modsFileName)

	file, err := os.ReadFile(name)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var mods []Mod
	err = json.Unmarshal(file, &mods)
	if err != nil {
		return nil, err
	}

	p.mods = mods
	return p.mods, nil
}

func (p *Profile) InstallMod(client *modrinth.Client, ctx context.Context, unsaved UnsavedMod) error {
	modsDir := filepath.Join(p.Directory, "/mods/")

	mods, err := p.Mods()
	if err != nil {
		return err
	}
	for _, mod := range mods {
		if mod.Id == unsaved.Id {
			return ErrModAlreadyInstalled
		}
	}

	err = os.Mkdir(modsDir, os.ModePerm)
	if err != nil {
		if !errors.Is(err, os.ErrExist) {
			return err
		}
	}

	tempFile, err := os.CreateTemp(modsDir, unsaved.File.Filename+".temp*")
	if err != nil {
		return err
	}
	tempPath := tempFile.Name()

	err = client.DownloadFile(ctx, unsaved.File, tempFile)
	if err != nil {
		tempFile.Close()
		os.Remove(tempPath)
		return err
	}

	tempFile.Close()
	err = os.Rename(tempPath, filepath.Join(modsDir, unsaved.File.Filename))
	if err != nil {
		tempFile.Close()
		return err
	}

	mods = append(mods, Mod{
		unsaved.Id,
		unsaved.VersionId,
		true,
		unsaved.Provider,
		unsaved.File.Filename,
	})

	p.mods = mods
	return nil
}

func (p *Profile) SaveMods() error {
	mccliDir := filepath.Join(p.Directory, mcCliDir)
	name := filepath.Join(mccliDir, modsFileName)

	err := os.Mkdir(mccliDir, os.ModePerm)
	if err != nil {
		if !errors.Is(err, os.ErrExist) {
			return err
		}
	}

	var f *os.File

	_, err = os.Stat(name)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			f, err = os.Create(name)
			if err != nil {
				return fmt.Errorf("failed to create %s: %w", name, err)
			}
		} else {
			return err
		}
	} else {
		f, err = os.OpenFile(name, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
		if err != nil {
			return fmt.Errorf("failed to open %s: %w", name, err)
		}
	}

	defer f.Close()

	mods, err := p.Mods()
	if err != nil {
		return err
	}

	if err := json.NewEncoder(f).Encode(mods); err != nil {
		return fmt.Errorf("failed to write json: %w", err)
	}

	return nil
}
