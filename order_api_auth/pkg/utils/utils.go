package utils

import (
	"crypto/rand"
	"encoding/hex"
	"math/big"
)

func GenerateSessionID() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func GetSMSCode() int {
	mx := big.NewInt(9000)
	n, _ := rand.Int(rand.Reader, mx)
	return int(n.Int64()) + 1000
}

func GetDevSMSCode() int {
	return 3245
}
