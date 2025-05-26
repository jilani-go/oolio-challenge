package models

// OrderItem represents a product with quantity in an order
type OrderItem struct {
	ProductID string  `json:"productId" bson:"productId"`
	Quantity  int     `json:"quantity" bson:"quantity"`
	Price     float64 `json:"price" bson:"price"` // Price at time of order
}

// Order represents a customer order
type Order struct {
	ID         string      `json:"id" bson:"_id,omitempty"`
	Items      []OrderItem `json:"items" bson:"items"`
	CouponCode string      `json:"couponCode,omitempty" bson:"couponCode,omitempty"`
}
