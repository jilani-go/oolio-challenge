package services

import (
	"github.com/jilani-go/glofox/internal/models"
	"github.com/jilani-go/glofox/internal/repository"
)

// ProductService defines the interface for product business logic
type ProductService interface {
	// GetAllProducts returns all available products
	GetAllProducts() ([]models.Product, error)
	
	// GetProductByID returns a product by its ID
	GetProductByID(id string) (*models.Product, error)
}

// ProductServiceImpl implements ProductService
type ProductServiceImpl struct {
	repo repository.ProductRepository
}

// NewProductService creates a new product service
func NewProductService(repo repository.ProductRepository) ProductService {
	return &ProductServiceImpl{
		repo: repo,
	}
}

// GetAllProducts returns all available products
func (s *ProductServiceImpl) GetAllProducts() ([]models.Product, error) {
	return s.repo.FindAll()
}

// GetProductByID returns a product by its ID
func (s *ProductServiceImpl) GetProductByID(id string) (*models.Product, error) {
	return s.repo.FindByID(id)
}
