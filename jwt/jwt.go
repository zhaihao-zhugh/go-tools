package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type TOKEN struct {
	UserID   string
	Username string
	Group    string
	Level    int
	jwt.StandardClaims
}

func ParseToken(tokenString string, sign []byte) (*TOKEN, error) {
	if tokenString == "" {
		return nil, errors.New("tokenString is empty")
	}

	token, err := jwt.ParseWithClaims(tokenString, &TOKEN{}, func(token *jwt.Token) (i interface{}, e error) {
		return sign, nil
	})
	if err != nil {
		return nil, err
	}

	if token != nil {
		if claims, ok := token.Claims.(*TOKEN); ok {
			return claims, nil
		}
		return nil, errors.New("token error")

	}
	return nil, errors.New("none token")
}

func CreateToken(t TOKEN, expires int64, sign []byte) (string, error) {
	t.StandardClaims = jwt.StandardClaims{
		NotBefore: time.Now().Unix() - 1000,
		ExpiresAt: time.Now().Unix() + expires,
		Issuer:    string(sign),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, t)
	return token.SignedString(sign)
}
