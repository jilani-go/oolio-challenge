package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jilani-go/glofox/internal/services"
)

// ProductHandler handles product-related requests
type ProductHandler struct {
	service services.ProductService
}

// NewProductHandler creates a new product handler
func NewProductHandler(service services.ProductService) *ProductHandler {
	return &ProductHandler{
		service: service,
	}
}

// ListProducts handles GET /api/product requests
// Returns a list of all available products
func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	// Get products from service
	modelProducts, err := h.service.GetAllProducts()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve products")
		return
	}

	// Convert model products to API products
	products := make([]Product, 0, len(modelProducts))
	for _, p := range modelProducts {
		products = append(products, Product{
			ID:       p.ID,
			Name:     p.Name,
			Price:    p.Price,
			Category: p.Category,
		})
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Encode and return the response
	if err := json.NewEncoder(w).Encode(products); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetProduct handles GET /api/product/{productId} requests
// Returns a single product by ID
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	// Extract the product ID from the URL path
	vars := mux.Vars(r)
	productID := vars["productId"]

	// Get product from service
	modelProduct, err := h.service.GetProductByID(productID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve product")
		return
	}

	// If product not found, return 404
	if modelProduct == nil {
		respondWithError(w, http.StatusNotFound, "Product not found")
		return
	}

	// Convert model product to API product
	product := Product{
		ID:       modelProduct.ID,
		Name:     modelProduct.Name,
		Price:    modelProduct.Price,
		Category: modelProduct.Category,
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Encode and return the response
	if err := json.NewEncoder(w).Encode(product); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// Helper function to respond with an error message
func respondWithError(w http.ResponseWriter, code int, message string) {
	response := ApiResponse{
		Code:    code,
		Type:    "error",
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}
