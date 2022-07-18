package util

import (
	"iot_server4/config"

	"github.com/golang-jwt/jwt"
)

func ParseToken(tokenString string) bool {
	tokenResult, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Conf.JwtSecret), nil
	})
	return tokenResult.Valid
}
