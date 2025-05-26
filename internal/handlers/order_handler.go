package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jilani-go/glofox/internal/models"
	"github.com/jilani-go/glofox/internal/services"
)

// OrderHandler handles order-related requests
type OrderHandler struct {
	orderService   services.OrderService
	productService services.ProductService
	validator      *validator.Validate
	promoService   services.PromoService
}

// NewOrderHandler creates a new order handler
func NewOrderHandler(orderService services.OrderService, productService services.ProductService, promoService services.PromoService) *OrderHandler {
	return &OrderHandler{
		orderService:   orderService,
		productService: productService,
		validator:      validator.New(),
		promoService:   promoService,
	}
}

// PlaceOrder handles POST /api/order requests
// Creates a new order with the provided items
func (h *OrderHandler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	// Check for API key (authentication)
	apiKey := r.Header.Get("api_key")
	if apiKey != "apitest" {
		respondWithError(w, http.StatusUnauthorized, "Invalid or missing API key")
		return
	}

	// Parse the request body into an OrderReq struct
	var orderReq OrderReq
	if err := json.NewDecoder(r.Body).Decode(&orderReq); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Validate the request
	if err := h.validator.Struct(orderReq); err != nil {
		respondWithError(w, http.StatusBadRequest, "Validation error: "+err.Error())
		return
	}

	// Convert request to domain model
	orderItems := make([]models.OrderItem, len(orderReq.Items))
	for i, item := range orderReq.Items {
		orderItems[i] = models.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
	}

	order := &models.Order{
		Items:      orderItems,
		CouponCode: orderReq.CouponCode,
	}
	valid, err := h.promoService.ValidatePromoCode(order.CouponCode)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error validating promo code: "+err.Error())
		return
	}
	if !valid {
		respondWithError(w, http.StatusBadRequest, "Invalid promo code")
		return
	}

	// Create the order via service
	createdOrder, err := h.orderService.CreateOrder(order)
	if err != nil {
		if errors.Is(err, services.ErrProductNotFound) {
			respondWithError(w, http.StatusBadRequest, "One or more products not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to create order: "+err.Error())
		return
	}

	// Get product details for response
	products := make([]Product, 0, len(createdOrder.Items))
	for _, item := range createdOrder.Items {
		product, err := h.productService.GetProductByID(item.ProductID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error retrieving product details")
			return
		}

		if product != nil {
			products = append(products, Product{
				ID:       product.ID,
				Name:     product.Name,
				Price:    product.Price,
				Category: product.Category,
			})
		}
	}

	// Create API response
	apiItems := make([]OrderItem, len(createdOrder.Items))
	for i, item := range createdOrder.Items {
		apiItems[i] = OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
	}

	// Create the order response
	orderResponse := Order{
		ID:       createdOrder.ID,
		Items:    apiItems,
		Products: products,
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// Encode and return the response
	if err := json.NewEncoder(w).Encode(orderResponse); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
