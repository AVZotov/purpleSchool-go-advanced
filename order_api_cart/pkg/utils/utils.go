package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/jinzhu/copier"
	pkgErr "order_api_cart/pkg/errors"
)

func GenerateSessionID() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("%w %v", pkgErr.ErrServiceUnavailable, err)
	}
	return hex.EncodeToString(bytes), nil
}

func GetFakeCode() int {
	return 3245
}

func ConvertToModel(dest, src interface{}) error {
	err := copier.Copy(dest, src)
	if err != nil {
		return err
	}
	return nil
}
