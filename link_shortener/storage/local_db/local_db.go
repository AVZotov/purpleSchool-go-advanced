package local_db

import (
	"io"
	"os"
)

const FILE = "verifyReq.json"

type FileHandler struct{}

func NewFileHandler() *FileHandler {
	return &FileHandler{}
}

func (h *FileHandler) Read() (io.ReadCloser, error) {
	return os.Open(FILE)
}

func (h *FileHandler) Write() (io.WriteCloser, error) {
	return os.OpenFile(FILE, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
}
