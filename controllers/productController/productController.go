package productController

import (
	"encoding/json"
	"github.com/fachryaa/project-assignment-synapsis-ecommerce/Responses"
	"github.com/fachryaa/project-assignment-synapsis-ecommerce/helper"
	"github.com/fachryaa/project-assignment-synapsis-ecommerce/models"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
)

func Create(w http.ResponseWriter, r *http.Request) {
	// mengambil inputan json
	var productInput models.Product
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&productInput); err != nil {
		response := Responses.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
		}
		helper.ResponseJson(w, response)

		return
	}
	defer r.Body.Close()

	// insert ke database
	if err := models.DB.Create(&productInput).Error; err != nil {
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
		Data:   productInput,
	}
	helper.ResponseJson(w, response)
}

func FindAll(w http.ResponseWriter, r *http.Request) {
	var result []models.Product
	if err := models.DB.Find(&result).Error; err != nil {
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

func FindById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productId := vars["productId"]

	var result models.Product

	if err := models.DB.First(&result, productId).Error; err != nil {
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

	response := Responses.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   result,
	}
	helper.ResponseJson(w, response)
}

func FindByCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productCategory := vars["category"]

	var result []models.Product

	if err := models.DB.Where("product_category LIKE ?", "%"+productCategory+"%").Find(&result).Error; err != nil {
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

	response := Responses.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   result,
	}
	helper.ResponseJson(w, response)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productId := vars["productId"]

	var product models.Product

	if err := models.DB.Delete(product, productId).Error; err != nil {
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

	response := Responses.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   product,
	}
	helper.ResponseJson(w, response)
}
