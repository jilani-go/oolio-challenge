package services

import (
	"errors"

	"github.com/jilani-go/glofox/internal/models"
	"github.com/jilani-go/glofox/internal/repository"
)

// Custom errors
var (
	ErrProductNotFound = errors.New("one or more products not found")
)

// OrderService defines the interface for order business logic
type OrderService interface {
	// CreateOrder validates and creates a new order
	CreateOrder(order *models.Order) (*models.Order, error)
	
	// ValidateOrderItems checks if all products in the order exist
	ValidateOrderItems(items []models.OrderItem) error
}

// OrderServiceImpl implements OrderService
type OrderServiceImpl struct {
	orderRepo   repository.OrderRepository
	productRepo repository.ProductRepository
}

// NewOrderService creates a new order service
func NewOrderService(orderRepo repository.OrderRepository, productRepo repository.ProductRepository) OrderService {
	return &OrderServiceImpl{
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

// CreateOrder validates and creates a new order
func (s *OrderServiceImpl) CreateOrder(order *models.Order) (*models.Order, error) {
	// First validate all order items exist
	if err := s.ValidateOrderItems(order.Items); err != nil {
		return nil, err
	}
	
	// Then create the order
	return s.orderRepo.Create(order)
}

// ValidateOrderItems checks if all products in the order exist
func (s *OrderServiceImpl) ValidateOrderItems(items []models.OrderItem) error {
	for _, item := range items {
		product, err := s.productRepo.FindByID(item.ProductID)
		if err != nil {
			return err
		}
		
		if product == nil {
			return ErrProductNotFound
		}
	}
	
	return nil
}
