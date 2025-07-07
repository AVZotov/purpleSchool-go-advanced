package local_storage

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

const TMPDIR = "tmp"

type FileHandler struct {
	WorkDir string
}

func newFileHandler(env string) (*FileHandler, error) {
	path, err := getFullPath(env)
	if err != nil {
		return nil, err
	}

	return &FileHandler{
		WorkDir: path,
	}, nil
}

func (h *FileHandler) load(name string) (io.ReadCloser, error) {
	file := filepath.Join(h.WorkDir, name)
	return os.Open(file)
}

func (h *FileHandler) save(name string) (io.WriteCloser, error) {
	file := filepath.Join(h.WorkDir, name)
	return os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
}

func (h *FileHandler) delete(name string) error {
	file := filepath.Join(h.WorkDir, name)
	return os.Remove(file)
}

func getFullPath(env string) (string, error) {
	switch env {
	case "dev":
		_, filename, _, _ := runtime.Caller(0)
		currentDir := filepath.Dir(filename)
		path := filepath.Join(currentDir, TMPDIR)
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return "", err
		}
		return path, nil
	case "prod":
		return os.TempDir(), nil
	default:
		return "", errors.New("unknown environment set")
	}
}
