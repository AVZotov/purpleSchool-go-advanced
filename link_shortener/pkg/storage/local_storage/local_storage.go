package local_storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"io"
	"log/slog"
	"os"
	"path"
	"strconv"
)

type Storage struct {
	FileHandler *FileHandler
	log         *slog.Logger
}

func NewStorage(devEnv string, log *slog.Logger) (*Storage, error) {
	log.With("link_shortener.pkg.storage.local_storage.local_storage.NewStorage()")
	fileHandler, err := newFileHandler(devEnv, log)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	storage := &Storage{
		FileHandler: fileHandler,
		log:         log,
	}
	log.Debug("new local storage created")

	return storage, nil
}

func (s *Storage) Save(email string, hash string) error {
	log := s.log.With("link_shortener.pkg.storage.local_storage.local_storage.Save()")
	fileName, err := getName(hash, log)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	if s.fileExists(fileName) {
		log.Warn("file already exists")
		return errors.New("file exists")
	}

	file, err := s.FileHandler.save(fileName)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	defer file.Close()

	bin := newBin(email, hash)

	payload, err := json.Marshal(bin)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	_, err = file.Write(payload)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	log.Debug("file saved to local storage")

	return nil
}

func (s *Storage) Load(hash string) (_ map[string]string, err error) {
	log := s.log.With("link_shortener.pkg.storage.local_storage.local_storage.Load()")
	details := make(map[string]string, 2)
	fileName, err := getName(hash, log)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	if !s.fileExists(fileName) {
		log.Error("file does not exist")
		return nil, errors.New("file does not exist")
	}

	file, err := s.FileHandler.load(fileName)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	defer func() {
		if err = file.Close(); err != nil {
			log.Error(err.Error(), "error closing file")
		}
	}()

	payload, err := io.ReadAll(file)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	var bin Bin
	err = json.Unmarshal(payload, &bin)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	details["email"] = bin.Email
	details["hash"] = bin.Hash

	log.Debug("file loaded from local storage")

	return details, nil
}

func (s *Storage) Delete(hash string) error {
	log := s.log.With("link_shortener.pkg.storage.local_storage.local_storage.Delete()")
	fileName, err := getName(hash, log)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	if !s.fileExists(fileName) {
		log.Error("file does not exist")
		return errors.New("file does not exist")
	}

	err = s.FileHandler.delete(fileName)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	log.Debug("file deleted from local storage")
	return nil
}

func getName(hash string, log *slog.Logger) (string, error) {
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
	log := s.log.With("link_shortener.pkg.storage.local_storage.local_storage.FileExists()")
	filePath := path.Join(s.FileHandler.WorkDir, fileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Debug(fmt.Sprintf("file %s does not exist", filePath))

		return false
	}
	log.Debug(fmt.Sprintf("file %s does exist", filePath))

	return true
}
