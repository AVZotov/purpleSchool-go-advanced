package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	pkgErrors "order_api_cart/pkg/errors"
)

func GenerateSessionID() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("%w %v", pkgErrors.ErrGeneratingSessionID, err)
	}
	return hex.EncodeToString(bytes), nil
}

func GetFakeCode() int {
	return 3245
}
