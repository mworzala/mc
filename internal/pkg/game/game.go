package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
)

type Manager interface {
}

var (
	procFileName = "proc.json"
)

type fileManager struct {
	Path    string             `json:"-"`
	Running map[string][]int64 `json:"running"`
}

func NewManager(dataDir string) (Manager, error) {
	procFile := path.Join(dataDir, procFileName)
	if _, err := os.Stat(procFile); errors.Is(err, fs.ErrNotExist) {
		return &fileManager{
			Path:    procFile,
			Running: make(map[string][]int64),
		}, nil
	}

	f, err := os.Open(procFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open proc file: %w", err)
	}
	defer f.Close()

	manager := fileManager{Path: procFile}
	if err := json.NewDecoder(f).Decode(&manager); err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", procFileName, err)
	}
	if manager.Running == nil {
		manager.Running = make(map[string][]int64)
	}
	return &manager, nil
}

//func (m *fileManager) updateProcesses() error {
//
//}

func (m *fileManager) Save() error {
	f, err := os.OpenFile(m.Path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("failed to open %s: %w", m.Path, err)
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(m); err != nil {
		return fmt.Errorf("failed to write json: %w", err)
	}

	return nil
}
