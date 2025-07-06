package local_storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"io"
	"os"
	"path"
	"strconv"
)

type Storage struct {
	FileHandler *FileHandler
}

func NewStorage(devEnv string) (*Storage, error) {
	fileHandler, err := newFileHandler(devEnv)
	if err != nil {
		return nil, err
	}

	return &Storage{
		FileHandler: fileHandler,
	}, nil
}

func (s *Storage) Save(email string, hash string) error {
	fileName, err := getName(email, hash)
	if err != nil {
		return err
	}

	if s.fileExists(fileName) {
		return errors.New("file already exists")
	}

	file, err := s.FileHandler.save(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	bin := newBin(email, hash)

	payload, err := json.Marshal(bin)
	if err != nil {
		return err
	}

	_, err = file.Write(payload)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) Load(email string, hash string) (_ *map[string]string, err error) {
	details := make(map[string]string, 2)
	fileName, err := getName(email, hash)
	if err != nil {
		return nil, err
	}

	if !s.fileExists(fileName) {
		return nil, errors.New("file does not exist")
	}

	file, err := s.FileHandler.load(fileName)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err = file.Close(); err != nil {
			fmt.Println("error closing file")
		}
		err = s.FileHandler.delete(fileName)
		if err != nil {
			fmt.Println("error deleting file")
			return
		}
	}()

	payload, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var bin Bin
	err = json.Unmarshal(payload, &bin)
	if err != nil {
		return nil, err
	}

	details["email"] = bin.Email
	details["hash"] = bin.Hash

	return &details, nil
}

func getName(email string, hash string) (string, error) {
	s := fmt.Sprintf("%s:%s", email, hash)
	hasher := fnv.New32a()
	_, err := hasher.Write([]byte(s))
	if err != nil {
		return "", err
	}
	name := strconv.Itoa(int(hasher.Sum32()))
	return fmt.Sprintf("%s.json", name), nil
}

func (s *Storage) fileExists(fileName string) bool {
	filePath := path.Join(s.FileHandler.WorkDir, fileName)
	_, err := os.Stat(filePath)
	if err == nil {
		return true
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return false
}
