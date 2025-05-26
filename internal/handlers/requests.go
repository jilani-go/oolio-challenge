package handlers

// Product represents a food product
type Product struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	ProductID string `json:"productId" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,min=1"`
}

// OrderReq represents the API request for placing an order
type OrderReq struct {
	CouponCode string      `json:"couponCode,omitempty"`
	Items      []OrderItem `json:"items" validate:"required,min=1,dive"`
}

// Order represents a placed order
type Order struct {
	ID       string      `json:"id"`
	Items    []OrderItem `json:"items"`
	Products []Product   `json:"products"`
}

// ApiResponse represents a general API response
type ApiResponse struct {
	Code    int    `json:"code"`
	Type    string `json:"type"`
	Message string `json:"message"`
}
