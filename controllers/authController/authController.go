package authController

import (
	"encoding/json"
	"github.com/fachryaa/project-assignment-synapsis-ecommerce/Responses"
	"github.com/fachryaa/project-assignment-synapsis-ecommerce/config"
	"github.com/fachryaa/project-assignment-synapsis-ecommerce/helper"
	"github.com/fachryaa/project-assignment-synapsis-ecommerce/models"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func Login(w http.ResponseWriter, r *http.Request) {
	// mengambil inputan json
	var userInput models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		response := Responses.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
		}
		helper.ResponseJson(w, response)

		return
	}
	defer r.Body.Close()

	// ambil data user berdasarkan username
	var user models.User
	if err := models.DB.Where("user_name = ?", userInput.UserName).First(&user).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			response := Responses.WebResponse{
				Code:   http.StatusNotFound,
				Status: "USERNAME NOT FOUND",
			}
			helper.ResponseJson(w, response)
			return
		default:
			response := Responses.WebResponse{
				Code:   http.StatusInternalServerError,
				Status: "INTERNAL ERROR",
			}
			helper.ResponseJson(w, response)
			return
		}
	}

	// cek apakah password valid
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password)); err != nil {
		response := Responses.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "USERNAME OR PASSWORD INVALID",
		}
		helper.ResponseJson(w, response)
		return
	}

	// proses pembuatan token jwt
	expTime := time.Now().Add(time.Minute * 60)
	claims := &config.JWTClaim{
		Id: user.UserId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "go-jwt-mux",
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}

	// deklarasi algoritma yang digunakan untuk signin
	tokenAlgo := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// signed token
	token, err := tokenAlgo.SignedString(config.JWT_KEY)
	if err != nil {
		response := Responses.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL ERROR",
		}
		helper.ResponseJson(w, response)
		return
	}

	// set token ke cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    token,
		HttpOnly: true,
	})

	response := Responses.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
	}
	helper.ResponseJson(w, response)
}

func Register(w http.ResponseWriter, r *http.Request) {
	// mengambil inputan json
	var userInput models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		response := Responses.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
		}
		helper.ResponseJson(w, response)

		return
	}
	defer r.Body.Close()

	// hash pass menggunakan bcrypt
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
	userInput.Password = string(hashPassword)

	// insert ke database
	if err := models.DB.Create(&userInput).Error; err != nil {
		response := Responses.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL ERROR",
		}
		helper.ResponseJson(w, response)

		return
	}

	response := Responses.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
	}
	helper.ResponseJson(w, response)

}

func Logout(w http.ResponseWriter, r *http.Request) {
	// hapus token yang ada di cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
	})

	response := Responses.WebResponse{
		Code:   http.StatusOK,
		Status: "LOGOUT SUCCESS",
	}
	helper.ResponseJson(w, response)

}
