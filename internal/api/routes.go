package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jilani-go/glofox/internal/handlers"
)

// SetupRoutes initializes the API routes
func SetupRoutes(productHandler *handlers.ProductHandler, orderHandler *handlers.OrderHandler) http.Handler {
	// Create router
	router := mux.NewRouter()

	// Product routes
	router.HandleFunc("/product", productHandler.ListProducts).Methods("GET")
	router.HandleFunc("/product/{productId}", productHandler.GetProduct).Methods("GET")

	// Order routes
	router.HandleFunc("/order", orderHandler.PlaceOrder).Methods("POST")

	return router
}
