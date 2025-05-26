package models

// OrderItem represents a product with quantity in an order
type OrderItem struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

// Order represents a customer order
type Order struct {
	ID         string      `json:"id"`
	Items      []OrderItem `json:"items"`
	CouponCode string      `json:"couponCode,omitempty"`
}
