package security

import (
	"crypto/sha256"
	"fmt"
	"time"
)

type Hash struct{}

func NewHashHandler() *Hash {
	return &Hash{}
}

func (h Hash) GetHash(email string) string {
	data := fmt.Sprintf("%s-%d", email, time.Now().UnixNano())
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(data)))
	return hash
}
