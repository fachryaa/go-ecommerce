package config

import "github.com/golang-jwt/jwt/v4"

var JWT_KEY = []byte("Vaguely-Instructive-Statement-08")

type JWTClaim struct {
	Id int64 `json:"id"`
	jwt.RegisteredClaims
}
