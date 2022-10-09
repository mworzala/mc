package util

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/k0kubun/go-ansi"
	"github.com/mworzala/mc-cli/internal/pkg/model"
	"github.com/schollz/progressbar/v3"
	"io"
	"net/http"
	"os"
	"path"
	"time"
)

type Context struct {
	RootDir string
	Ctx     context.Context
}

func NewContext(dir string) *Context {
	return &Context{
		dir,
		context.Background(),
	}
}

func (c Context) CreateProgressBar(size int64) *progressbar.ProgressBar {
	if size < 1 {
		size = -1
	}
	return progressbar.NewOptions64(size,
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetDescription(" "),
		progressbar.OptionShowElapsedTimeOnFinish(),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]━[reset]",
			SaucerHead:    "[yellow]━[reset]",
			SaucerPadding: " ",
		}),
	)
}

func (c Context) ReadOrDownload(id, path string, file model.Download, ptr interface{}) error {
	if _, err := os.Stat(path); err == nil {
		// Download if not exists
		return c.Read(path, ptr)
	} else if errors.Is(err, os.ErrNotExist) {
		// Download the file with a progress bar and buffer for json parsing
		data := new(bytes.Buffer)
		bar := c.CreateProgressBar(file.Size)

		if id != "" {
			fmt.Println("\n" + id)
		}
		err = c.DownloadFile(path, file, bar, data)
		if err != nil {
			return err
		}

		// Newline after progress bar
		fmt.Println("")

		// Parse json from downloaded data
		return json.NewDecoder(data).Decode(ptr)
	} else {
		return err
	}
}

func (c Context) Download(id, file string, dl model.Download) error {
	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
		// Download the file with a progress bar
		bar := c.CreateProgressBar(dl.Size)
		if id != "" {
			fmt.Println(id)
		}
		err = c.DownloadFile(file, dl, bar)
		if err != nil {
			return err
		}

		// Newline after progress bar
		fmt.Println("")
	}

	return nil
}

func (c Context) DownloadFile(file string, dl model.Download, listeners ...io.Writer) error {
	// Create parent directory
	err := os.MkdirAll(path.Dir(file), 0755)
	if err != nil {
		return err
	}

	// Open request
	res, err := http.Get(dl.Url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Create file
	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// Copy data to file, hash, and listeners
	hash := sha1.New()
	writers := listeners
	writers = append(writers, f)
	if dl.Sha1 != "" {
		writers = append(writers, hash)
	}
	_, err = io.Copy(io.MultiWriter(writers...), res.Body)
	if err != nil {
		return err
	}

	// Validate hash if present
	if dl.Sha1 != "" {
		h := fmt.Sprintf("%x", hash.Sum(nil))
		if h != dl.Sha1 {
			return fmt.Errorf("download hash mismatch: %s != %s", h, dl.Sha1)
		}
	}

	return nil
}

func (c Context) Read(path string, ptr interface{}) error {
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	return json.NewDecoder(file).Decode(ptr)
}

func (c Context) Deadline() (deadline time.Time, ok bool) {
	return c.Ctx.Deadline()
}

func (c Context) Done() <-chan struct{} {
	return c.Ctx.Done()
}

func (c Context) Err() error {
	return c.Ctx.Err()
}

func (c Context) Value(key any) any {
	return c.Ctx.Value(key)
}
