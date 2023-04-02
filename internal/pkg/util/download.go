package util

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
)

type FileDownload struct {
	Sha1 string `json:"sha1"`
	Size int64  `json:"size"`
	Url  string `json:"url"`
}

func ReadOrDownload(id, file string, dl FileDownload, ptr interface{}) error {
	if _, err := os.Stat(file); err == nil {
		return ReadFile(file, ptr)
	} else if errors.Is(err, fs.ErrNotExist) {
		data := new(bytes.Buffer)

		println("download", dl.Url) //todo remove me/add some callback for progress
		if err := downloadFile(file, dl, data); err != nil {
			return err
		}

		return json.NewDecoder(data).Decode(ptr)
	} else {
		return err
	}
}

func Download(id, file string, dl FileDownload) error {
	if _, err := os.Stat(file); errors.Is(err, fs.ErrNotExist) {
		println("download", dl.Url) //todo remove me/add some callback for progress
		return downloadFile(file, dl)
	}
	return nil
}

func downloadFile(file string, download FileDownload, listeners ...io.Writer) error {
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

			return fmt.Errorf("FileDownload hash mismatch: %s != %s", h, download.Sha1)
		}
	}

	return nil
}

func ReadFile(file string, ptr interface{}) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(ptr)
}
