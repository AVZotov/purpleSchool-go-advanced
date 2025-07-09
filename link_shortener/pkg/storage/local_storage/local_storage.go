package local_storage

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io"
	t "link_shortener/pkg/storage/types"
	"os"
	"path"
	"strconv"
)

type Storage struct {
	FileHandler *Handler
	Log         t.Logger
}

func New(devEnv string, log t.Logger) (*Storage, error) {
	const fn = "pkg.storage.local_storage.local_storage.new"
	s := &Storage{
		Log: log,
	}

	s.Log.With(fn)

	fileHandler, err := newHandler(devEnv, s.Log)
	if err != nil {
		s.Log.Error(err.Error())
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	s.FileHandler = fileHandler

	log.Debug("new local storage created")

	return s, nil
}

func (s *Storage) Save(email string, hash string) error {
	const fn = "pkg.storage.local_storage.local_storage.save"
	s.Log.With(fn)
	fileName, err := getName(hash, s.Log)
	if err != nil {
		s.Log.Error(fmt.Sprintf("%s: %s", fn, err.Error()))
		return fmt.Errorf("%s: %w", fn, err)
	}

	if s.fileExists(fileName) {
		s.Log.Warn(fmt.Sprintf("%s:%s already exists", fn, fileName))
		return fmt.Errorf("%s:%s already exists", fn, fileName)
	}

	file, err := s.FileHandler.save(fileName)
	if err != nil {
		s.Log.Error(fmt.Sprintf("%s: %s", fn, err.Error()))
		return fmt.Errorf("%s: %w", fn, err)
	}
	defer file.Close()

	bin := newBin(email, hash)

	payload, err := json.Marshal(bin)
	if err != nil {
		s.Log.Error(fmt.Sprintf("%s: %s", fn, err.Error()))
		return fmt.Errorf("%s: %w", fn, err)
	}

	_, err = file.Write(payload)
	if err != nil {
		s.Log.Error(fmt.Sprintf("%s: %s", fn, err.Error()))
		return fmt.Errorf("%s: %w", fn, err)
	}

	s.Log.Debug("file saved to local storage")

	return nil
}

func (s *Storage) Load(hash string) (map[string]string, error) {
	const fn = "pkg.storage.local_storage.local_storage.load"
	s.Log.With(fn)
	details := make(map[string]string, 2)
	fileName, err := getName(hash, s.Log)
	if err != nil {
		s.Log.Error(fmt.Sprintf("%s: %s", fn, err.Error()))
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	if !s.fileExists(fileName) {
		s.Log.Warn(fmt.Sprintf("%s:%s does not exists", fn, fileName))
		return nil, fmt.Errorf("%s:%s does not exists", fn, fileName)
	}

	file, err := s.FileHandler.load(fileName)
	if err != nil {
		s.Log.Error(fmt.Sprintf("%s: %s", fn, err.Error()))
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	defer func() {
		if err = file.Close(); err != nil {
			s.Log.Error(fn, err.Error(), "error closing file")
		}
	}()

	payload, err := io.ReadAll(file)
	if err != nil {
		s.Log.Error(fmt.Sprintf("%s: %s", fn, err.Error()))
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	var bin Bin
	err = json.Unmarshal(payload, &bin)
	if err != nil {
		s.Log.Error(fmt.Sprintf("%s: %s", fn, err.Error()))
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	details["email"] = bin.Email
	details["hash"] = bin.Hash

	s.Log.Debug(fn, "file loaded from local storage")

	return details, nil
}

func (s *Storage) Delete(hash string) error {
	const fn = "link_shortener.pkg.storage.local_storage.local_storage.Delete"
	s.Log.With(fn)
	fileName, err := getName(hash, s.Log)
	if err != nil {
		s.Log.Error(fmt.Sprintf("%s: %s", fn, err.Error()))
		return fmt.Errorf("%s: %w", fn, err)
	}

	if !s.fileExists(fileName) {
		s.Log.Warn(fmt.Sprintf("%s:%s does not exists", fn, fileName))
		return fmt.Errorf("%s:%s does not exists", fn, fileName)
	}

	err = s.FileHandler.delete(fileName)
	if err != nil {
		s.Log.Error(fmt.Sprintf("%s: %s", fn, err.Error()))
		return fmt.Errorf("%s: %w", fn, err)
	}

	s.Log.Debug("file deleted from local storage")
	return nil
}

func getName(hash string, log t.Logger) (string, error) {
	log.With("link_shortener.pkg.storage.local_storage.local_storage.getName()")
	hasher := fnv.New32a()
	_, err := hasher.Write([]byte(hash))
	if err != nil {
		log.Error(err.Error())
		return "", err
	}
	name := strconv.Itoa(int(hasher.Sum32()))

	log.Debug(fmt.Sprintf("file name generated: %s.json", name))

	return fmt.Sprintf("%s.json", name), nil
}

func (s *Storage) fileExists(fileName string) bool {
	const fn = "link_shortener.pkg.storage.local_storage.local_storage.fileExists"
	s.Log.With(fn)
	filePath := path.Join(s.FileHandler.WorkDir, fileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		s.Log.Debug(fmt.Sprintf("%s:file %s does not exist", fn, filePath))

		return false
	}
	s.Log.Debug(fmt.Sprintf("%s:file %s exists", fn, filePath))

	return true
}
