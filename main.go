package main

import (
	"github.com/fachryaa/project-assignment-synapsis-ecommerce/controllers/authController"
	"github.com/fachryaa/project-assignment-synapsis-ecommerce/controllers/cartController"
	"github.com/fachryaa/project-assignment-synapsis-ecommerce/controllers/orderController"
	"github.com/fachryaa/project-assignment-synapsis-ecommerce/controllers/productController"
	"github.com/fachryaa/project-assignment-synapsis-ecommerce/middlewares"
	"github.com/fachryaa/project-assignment-synapsis-ecommerce/models"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	models.ConnectDatabase()

	r := mux.NewRouter()

	r.HandleFunc("/login", authController.Login).Methods("POST")
	r.HandleFunc("/register", authController.Register).Methods("POST")
	r.HandleFunc("/logout", authController.Logout).Methods("GET")

	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/products", productController.Create).Methods("POST")
	api.HandleFunc("/products", productController.FindByCategory).Queries("category", "{category}").Methods("GET")
	api.HandleFunc("/products", productController.FindAll).Methods("GET")
	api.HandleFunc("/products/{productId}", productController.FindById).Methods("GET")
	api.HandleFunc("/products/{productId}", productController.Delete).Methods("GET")

	api.HandleFunc("/cart", cartController.Create).Methods("POST")
	api.HandleFunc("/cart", cartController.FindAll).Methods("GET")
	api.HandleFunc("/cart/{cartId}", cartController.UpdateAmount).Methods("PUT")
	api.HandleFunc("/cart/{cartId}", cartController.DeleteCart).Methods("DELETE")

	api.HandleFunc("/order", orderController.FindAll).Methods("GET")
	api.HandleFunc("/order/checkout", orderController.CheckoutOrder).Methods("GET")

	api.Use(middlewares.JWTMiddleware)

	log.Fatal(http.ListenAndServe("localhost:3000", r))
}
