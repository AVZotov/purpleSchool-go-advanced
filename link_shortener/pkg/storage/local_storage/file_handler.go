package local_storage

import (
	"fmt"
	"io"
	t "link_shortener/pkg/storage/types"
	"os"
	"path/filepath"
	"runtime"
)

const TMPDIR = "tmp"

type Handler struct {
	WorkDir string
	Log     t.Logger
}

func newHandler(env string, log t.Logger) (*Handler, error) {
	const fn = "pkg.storage.local_storage.file_handler.newHandler"
	fh := &Handler{
		Log: log,
	}

	fh.Log.With(fn)
	path, err := getFullPath(env, log)
	if err != nil {
		log.Error(err.Error())
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	fh.WorkDir = path

	log.Debug("fileHandler initialized")

	return fh, nil
}

func (h *Handler) load(name string) (io.ReadCloser, error) {
	const fn = "pkg.storage.local_storage.file_handler.load"
	h.Log.With(fn)

	filePath := filepath.Join(h.WorkDir, name)

	file, err := os.Open(filePath)
	if err != nil {
		h.Log.Error(err.Error())
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	h.Log.Debug("file opened for reading")

	return file, nil
}

func (h *Handler) save(name string) (io.WriteCloser, error) {
	const fn = "pkg.storage.local_storage.file_handler.save"
	h.Log.With(fn)
	filePath := filepath.Join(h.WorkDir, name)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		h.Log.Error(err.Error())
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	h.Log.Debug("file opened for writing")

	return file, nil
}

func (h *Handler) delete(name string) error {
	const fn = "pkg.storage.local_storage.file_handler.delete"
	h.Log.With(fn)
	file := filepath.Join(h.WorkDir, name)
	if err := os.Remove(file); err != nil {
		h.Log.Error(err.Error())
		return fmt.Errorf("%s: %w", fn, err)
	}

	h.Log.Debug("file deleted")

	return nil
}

func getFullPath(env string, log t.Logger) (string, error) {
	const fn = "pkg.storage.local_storage.file_handler.getFullPath"
	log.With(fn)

	switch env {
	case "dev":
		_, filename, _, _ := runtime.Caller(0)
		currentDir := filepath.Dir(filename)
		path := filepath.Join(currentDir, TMPDIR)
		err := os.MkdirAll(path, 0755)
		if err != nil {
			log.Error(err.Error())
			return "", fmt.Errorf("%s: %w", fn, err)
		}
		return path, nil
	case "prod":
		return os.TempDir(), nil
	default:
		log.Error(fmt.Sprintf("unknown env type: %s", env))
		return "", fmt.Errorf("%s: %w", fn, fmt.Errorf(
			"unknown env type: %s", env))
	}
}
