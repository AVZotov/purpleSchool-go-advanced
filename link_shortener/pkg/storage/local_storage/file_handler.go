package local_storage

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
)

const TMPDIR = "tmp"

type FileHandler struct {
	WorkDir string
	log     *slog.Logger
}

func newFileHandler(env string, log *slog.Logger) (*FileHandler, error) {
	log.With("link_shortener.pkg.storage.local_storage.file_handler.newFileHandler()")
	path, err := getFullPath(env, log)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	fh := &FileHandler{
		WorkDir: path,
		log:     log,
	}

	log.Debug("file handler initialized")

	return fh, nil
}

func (h *FileHandler) load(name string) (io.ReadCloser, error) {
	log := h.log.With("link_shortener.pkg.storage.local_storage.file_handler.load()")
	filePath := filepath.Join(h.WorkDir, name)

	file, err := os.Open(filePath)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	log.Debug("file opened for reading")

	return file, nil
}

func (h *FileHandler) save(name string) (io.WriteCloser, error) {
	log := h.log.With("link_shortener.pkg.storage.local_storage.file_handler.save")
	filePath := filepath.Join(h.WorkDir, name)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	log.Debug("file opened for writing")

	return file, nil
}

func (h *FileHandler) delete(name string) error {
	log := h.log.With("link_shortener.pkg.storage.local_storage.file_handler.delete()")
	file := filepath.Join(h.WorkDir, name)
	if err := os.Remove(file); err != nil {
		log.Error(err.Error())
		return err
	}

	log.Debug("file deleted")

	return nil
}

func getFullPath(env string, log *slog.Logger) (string, error) {
	log.With("link_shortener.pkg.storage.local_storage.file_handler.getFullPath()")
	switch env {
	case "dev":
		_, filename, _, _ := runtime.Caller(0)
		currentDir := filepath.Dir(filename)
		path := filepath.Join(currentDir, TMPDIR)
		err := os.MkdirAll(path, 0755)
		if err != nil {
			log.Error(err.Error())
			return "", err
		}
		return path, nil
	case "prod":
		return os.TempDir(), nil
	default:
		log.Error(fmt.Sprintf("unknown env type: %s", env))
		return "", errors.New("unknown environment set")
	}
}
