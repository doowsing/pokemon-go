package utils

import (
	"github.com/dgrijalva/jwt-go"
	"pokemon/pkg/config"
	"time"
)

var jwtSecret = []byte(config.Config().SecretKey)

type Claims struct {
	ID      int    `json:"id"`
	Account string `json:"account"`
	jwt.StandardClaims
}

func GenerateToken(ID int, Account string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour).Unix()

	claims := Claims{
		ID,
		Account,
		jwt.StandardClaims{
			ExpiresAt: expireTime,
			Issuer:    "poke",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
