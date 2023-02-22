package helper

import (
	"github.com/fachryaa/project-assignment-synapsis-ecommerce/config"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
)

func GetUserInfo(r *http.Request) (config.JWTClaim, error) {
	c, err := r.Cookie("token")
	// mengambil token value
	tokenString := c.Value

	claims := &config.JWTClaim{}

	// parsing token jwt
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return config.JWT_KEY, nil
	})

	if claims, ok := token.Claims.(*config.JWTClaim); ok && token.Valid {
		return *claims, nil
	} else {
		return *claims, err
	}
}
