package cartController

import (
	"encoding/json"
	"fmt"
	"github.com/fachryaa/project-assignment-synapsis-ecommerce/Responses"
	"github.com/fachryaa/project-assignment-synapsis-ecommerce/helper"
	"github.com/fachryaa/project-assignment-synapsis-ecommerce/models"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
)

func Create(w http.ResponseWriter, r *http.Request) {
	// get user info from token
	userInfo, err := helper.GetUserInfo(r)
	if err != nil {
		response := Responses.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: "UNAUTHORIZED",
		}
		helper.ResponseJson(w, response)
		return
	}
	userId := userInfo.Id

	// mengambil inputan json
	var cartInput models.Cart
	cartInput.UserId = userId // assign value
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&cartInput); err != nil {
		response := Responses.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
		}
		helper.ResponseJson(w, response)

		return
	}
	defer r.Body.Close()

	// get detail product
	var product models.Product
	if err := models.DB.First(&product, cartInput.ProductId).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			response := Responses.WebResponse{
				Code:   http.StatusNotFound,
				Status: "Data tidak ditemukan",
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

	// cek product in cart
	var productInCart models.Cart
	err = models.DB.Where("user_id = ? AND product_id = ?", userId, cartInput.ProductId).First(&productInCart).Error
	fmt.Println(productInCart)

	// cek if stock ready
	productStock := product.ProductStock
	amount := cartInput.Amount + productInCart.Amount
	if amount > productStock {
		response := Responses.WebResponse{
			Code:   http.StatusNotAcceptable,
			Status: "STOCK IS NOT READY",
		}
		helper.ResponseJson(w, response)
		return
	}

	// if product is not in cart
	if err == gorm.ErrRecordNotFound {
		// insert database
		cartInput.TotalPrice = amount * product.ProductPrice
		if err := models.DB.Create(&cartInput).Error; err != nil {
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
			Data:   cartInput,
		}
		helper.ResponseJson(w, response)
	} else {
		// if product is already in cart
		productInCart.Amount = amount
		productInCart.TotalPrice = amount * product.ProductPrice
		if err := models.DB.Save(&productInCart).Error; err != nil {
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
			Data:   productInCart,
		}
		helper.ResponseJson(w, response)
	}
}

func FindAll(w http.ResponseWriter, r *http.Request) {
	// get user info from token
	userInfo, err := helper.GetUserInfo(r)
	if err != nil {
		response := Responses.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: "UNAUTHORIZED",
		}
		helper.ResponseJson(w, response)
		return
	}
	userId := userInfo.Id

	// get carts
	var result []models.Cart
	if err := models.DB.Where("user_id = ?", userId).Find(&result).Error; err != nil {
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
		Data:   result,
	}
	helper.ResponseJson(w, response)
}

func UpdateAmount(w http.ResponseWriter, r *http.Request) {
	// get user info from token
	userInfo, err := helper.GetUserInfo(r)
	if err != nil {
		response := Responses.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: "UNAUTHORIZED",
		}
		helper.ResponseJson(w, response)
		return
	}
	userId := userInfo.Id

	// get body json
	var body models.Cart
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		response := Responses.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
		}
		helper.ResponseJson(w, response)

		return
	}
	inputAmount := body.Amount
	defer r.Body.Close()

	// get params
	vars := mux.Vars(r)
	cartId := vars["cartId"]

	// get cart by cartId
	var cart models.Cart
	if err := models.DB.Where("cart_id = ? AND user_id = ?", cartId, userId).First(&cart).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			response := Responses.WebResponse{
				Code:   http.StatusNotFound,
				Status: "Data tidak ditemukan",
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

	// get detail product
	var product models.Product
	if err := models.DB.First(&product, cart.ProductId).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			response := Responses.WebResponse{
				Code:   http.StatusNotFound,
				Status: "Data tidak ditemukan",
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
	productStock := product.ProductStock

	// check stock is ready
	if inputAmount > productStock {
		response := Responses.WebResponse{
			Code:   http.StatusNotAcceptable,
			Status: "STOCK IS NOT READY",
		}
		helper.ResponseJson(w, response)
		return
	}

	// update cart if stock is ready
	cart.Amount = inputAmount
	cart.TotalPrice = cart.Amount * product.ProductPrice
	if err := models.DB.Save(&cart).Error; err != nil {
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
		Data:   cart,
	}
	helper.ResponseJson(w, response)
}

func DeleteCart(w http.ResponseWriter, r *http.Request) {
	// get user info from token
	userInfo, err := helper.GetUserInfo(r)
	if err != nil {
		response := Responses.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: "UNAUTHORIZED",
		}
		helper.ResponseJson(w, response)
		return
	}
	userId := userInfo.Id

	// get params
	vars := mux.Vars(r)
	cartId := vars["cartId"]

	// get cart by cartId
	var cart models.Cart
	if err := models.DB.Where("cart_id = ? AND user_id = ?", cartId, userId).First(&cart).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			response := Responses.WebResponse{
				Code:   http.StatusNotFound,
				Status: "Data tidak ditemukan",
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

	// delete cart
	if err := models.DB.Delete(cart).Error; err != nil {
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
		Data:   cart,
	}
	helper.ResponseJson(w, response)
}
