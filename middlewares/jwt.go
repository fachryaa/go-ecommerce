package middlewares

import (
	"github.com/fachryaa/project-assignment-synapsis-ecommerce/Responses"
	"github.com/fachryaa/project-assignment-synapsis-ecommerce/config"
	"github.com/fachryaa/project-assignment-synapsis-ecommerce/helper"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				response := Responses.WebResponse{
					Code:   http.StatusUnauthorized,
					Status: "UNAUTHORIZED",
				}
				helper.ResponseJson(w, response)
				return
			}
		}

		// mengambil token value
		tokenString := c.Value

		claims := &config.JWTClaim{}

		// parsing token jwt
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return config.JWT_KEY, nil
		})

		if err != nil {
			v, _ := err.(*jwt.ValidationError)
			switch v.Errors {
			case jwt.ValidationErrorSignatureInvalid:
				// token invalid
				response := Responses.WebResponse{
					Code:   http.StatusUnauthorized,
					Status: "UNAUTHORIZED",
				}
				helper.ResponseJson(w, response)
				return
			case jwt.ValidationErrorExpired:
				// token expired
				response := Responses.WebResponse{
					Code:   http.StatusUnauthorized,
					Status: "UNAUTHORIZED",
				}
				helper.ResponseJson(w, response)
				return
			default:
				response := Responses.WebResponse{
					Code:   http.StatusUnauthorized,
					Status: "UNAUTHORIZED",
				}
				helper.ResponseJson(w, response)
				return
			}
		}

		if !token.Valid {
			response := Responses.WebResponse{
				Code:   http.StatusUnauthorized,
				Status: "UNAUTHORIZED",
			}
			helper.ResponseJson(w, response)
			return
		}

		next.ServeHTTP(w, r)
	})
}
