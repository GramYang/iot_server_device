package util

import (
	"iot_server8/config"
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateToken() string {
	secret := []byte(config.Conf.JwtSecret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"nbf": time.Now().Unix(),
	})
	tokenString, _ := token.SignedString(secret)
	return tokenString
}
