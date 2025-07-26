package jwt

import "github.com/golang-jwt/jwt/v5"

func Create(secret, phone string) (string, error) {
	jwtSecret := []byte(secret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"phone": phone,
	})
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ParsePhone(token, secret string) (bool, string) {
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return false, ""
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)["phone"]
	if !ok {
		return false, ""
	}
	phone, ok := claims.(string)
	if !ok {
		return false, ""
	}

	return jwtToken.Valid, phone
}
