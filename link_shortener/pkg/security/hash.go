package security

import (
	"crypto/md5"
	"fmt"
	"time"
)

type Hash struct{}

func NewHash() *Hash {
	return &Hash{}
}

func (h Hash) GetHash(email string) string {
	data := fmt.Sprintf("%s-%d", email, time.Now().Unix())
	hash := fmt.Sprintf("%x", md5.Sum([]byte(data)))
	return hash
}
