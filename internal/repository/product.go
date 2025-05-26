package repository

import (
	"sync"

	"github.com/jilani-go/glofox/internal/models"
)

// InMemoryProductRepository implements ProductRepository using in-memory storage
type InMemoryProductRepository struct {
	products []models.Product
	mutex    sync.RWMutex // Add RWMutex for thread safety
}

// NewInMemoryProductRepository creates a new repository with predefined products
func NewInMemoryProductRepository() *InMemoryProductRepository {
	// Initialize with some products
	products := []models.Product{
		{
			ID:       "1",
			Name:     "Waffle with Berries",
			Price:    6.5,
			Category: "Waffle",
		},
		{
			ID:       "2",
			Name:     "Vanilla Bean Crème Brûlée",
			Price:    7.0,
			Category: "Crème Brûlée",
		},
		{
			ID:       "3",
			Name:     "Macaron Mix of Five",
			Price:    8.0,
			Category: "Macaron",
		},
		{
			ID:       "4",
			Name:     "Classic Tiramisu",
			Price:    5.5,
			Category: "Tiramisu",
		},
		{
			ID:       "5",
			Name:     "Pistachio Baklava",
			Price:    4.0,
			Category: "Baklava",
		},
		{
			ID:       "6",
			Name:     "Lemon Meringue Pie",
			Price:    5.0,
			Category: "Pie",
		},
		{
			ID:       "7",
			Name:     "Red Velvet Cake",
			Price:    4.5,
			Category: "Cake",
		},
		{
			ID:       "8",
			Name:     "Salted Caramel Brownie",
			Price:    4.5,
			Category: "Brownie",
		},
		{
			ID:       "9",
			Name:     "Vanilla Panna Cotta",
			Price:    6.5,
			Category: "Panna Cotta",
		},
	}

	return &InMemoryProductRepository{
		products: products,
	}
}

// FindAll returns all products
func (r *InMemoryProductRepository) FindAll() ([]models.Product, error) {
	r.mutex.RLock()         // Use read lock for read-only operations
	defer r.mutex.RUnlock() // Ensure unlock happens even if there's a panic

	// Create a copy of the slice to prevent data races
	productsCopy := make([]models.Product, len(r.products))
	copy(productsCopy, r.products)

	return productsCopy, nil
}

// FindByID returns a product by its ID
func (r *InMemoryProductRepository) FindByID(id string) (*models.Product, error) {
	r.mutex.RLock()         // Use read lock for read-only operations
	defer r.mutex.RUnlock() // Ensure unlock happens even if there's a panic

	for _, product := range r.products {
		if product.ID == id {
			// Create a copy to prevent data races
			productCopy := product
			return &productCopy, nil
		}
	}
	return nil, nil // Not found
}
