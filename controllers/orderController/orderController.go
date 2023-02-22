package orderController

import (
	"fmt"
	"github.com/fachryaa/project-assignment-synapsis-ecommerce/Responses"
	"github.com/fachryaa/project-assignment-synapsis-ecommerce/helper"
	"github.com/fachryaa/project-assignment-synapsis-ecommerce/models"
	"net/http"
)

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

	// get orders
	var orders []models.Order
	if err := models.DB.Where("user_id = ?", userId).Find(&orders).Error; err != nil {
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
		Data:   orders,
	}
	helper.ResponseJson(w, response)
}

func CheckoutOrder(w http.ResponseWriter, r *http.Request) {
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

	// get user info from db
	var user models.User
	if err := models.DB.First(&user, userId).Error; err != nil {
		response := Responses.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL ERROR",
		}
		helper.ResponseJson(w, response)

		return
	}

	// get cart
	var carts []models.Cart
	rowsCart, err := models.DB.Where("user_id = ?", userId).Find(&carts).Rows()
	helper.PanicIfError(err)
	defer rowsCart.Close()

	// calculate for total price
	var totalPrice uint64
	for _, v := range carts {
		totalPrice += v.TotalPrice
	}

	// make new order
	var newOrder models.Order
	newOrder.FullName = user.FullName
	newOrder.UserId = user.UserId
	newOrder.TotalPrice = totalPrice
	newOrder.Address = user.Address
	newOrder.Phone = user.Phone

	// insert database
	if err := models.DB.Create(&newOrder).Error; err != nil {
		response := Responses.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL ERROR",
			Data:   err.Error(),
		}
		helper.ResponseJson(w, response)

		return
	}

	// get product detail
	var productIds []int64
	for _, v := range carts {
		productIds = append(productIds, v.ProductId)
	}

	// reduce stock product
	var products []models.Product
	rows, err := models.DB.Find(&products, productIds).Rows()
	helper.PanicIfError(err)
	defer rows.Close()

	i := 0
	for rows.Next() {
		var product models.Product
		models.DB.ScanRows(rows, &product)

		// update stock
		product.ProductStock -= carts[i].Amount
		fmt.Println(product, carts[i])
		i++
		if err := models.DB.Save(product).Error; err != nil {
			response := Responses.WebResponse{
				Code:   http.StatusInternalServerError,
				Status: "INTERNAL ERROR",
				Data:   err.Error(),
			}
			helper.ResponseJson(w, response)

			return
		}
	}

	// delete cart
	for rowsCart.Next() {
		var cart models.Cart
		models.DB.ScanRows(rowsCart, &cart)
		fmt.Println(cart)

		if err := models.DB.Delete(cart).Error; err != nil {
			response := Responses.WebResponse{
				Code:   http.StatusInternalServerError,
				Status: "INTERNAL ERROR",
				Data:   err.Error(),
			}
			helper.ResponseJson(w, response)

			return
		}
	}

	response := Responses.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   newOrder,
	}
	helper.ResponseJson(w, response)
}
