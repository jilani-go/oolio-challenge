package repository

import (
	"github.com/jilani-go/glofox/internal/models"
)

// ProductRepository defines the interface for product data operations
type ProductRepository interface {
	FindAll() ([]models.Product, error)
	FindByID(id string) (*models.Product, error)
}

// OrderRepository defines the interface for order data operations
type OrderRepository interface {
	Create(order *models.Order) (*models.Order, error)
}
