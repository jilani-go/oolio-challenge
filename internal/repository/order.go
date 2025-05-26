package repository

import (
	"sync"

	"github.com/google/uuid"
	"github.com/jilani-go/glofox/internal/models"
)

// InMemoryOrderRepository implements OrderRepository using in-memory storage
type InMemoryOrderRepository struct {
	orders      []models.Order
	productRepo ProductRepository
	mutex       sync.RWMutex // Add mutex for thread safety
}

// NewInMemoryOrderRepository creates a new repository for orders
func NewInMemoryOrderRepository(productRepo ProductRepository) *InMemoryOrderRepository {
	return &InMemoryOrderRepository{
		orders:      []models.Order{},
		productRepo: productRepo,
	}
}

// Create adds a new order
func (r *InMemoryOrderRepository) Create(order *models.Order) (*models.Order, error) {
	r.mutex.Lock()         // Lock for writing
	defer r.mutex.Unlock() // Ensure unlock happens even if there's a panic

	// Generate a new UUID for the order
	order.ID = uuid.New().String()

	// Add to storage
	r.orders = append(r.orders, *order)

	return order, nil
}
