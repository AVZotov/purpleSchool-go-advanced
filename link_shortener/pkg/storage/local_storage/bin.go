package local_storage

import "strings"

type Bin struct {
	Email string `json:"email" validator:"required,email"`
	Hash  string `json:"hash" validator:"required"`
}

func newBin(email string, hash string) *Bin {
	return &Bin{
		Email: strings.ToLower(email),
		Hash:  hash,
	}
}
