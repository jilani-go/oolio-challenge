package services

import (
	"time"

	"github.com/jilani-go/glofox/internal/models"
	"github.com/jilani-go/glofox/internal/repository"
)

// ClassService defines operations for class business logic
type ClassService interface {
	CreateClass(name string, startDate, endDate time.Time, capacity int) (models.Class, error)
}

// ClassServiceImpl implements ClassService
type ClassServiceImpl struct {
	classRepo   repository.ClassRepository
	bookingRepo repository.BookingRepository
}

// NewClassService creates a new class service
func NewClassService(classRepo repository.ClassRepository, bookingRepo repository.BookingRepository) ClassService {
	return &ClassServiceImpl{
		classRepo:   classRepo,
		bookingRepo: bookingRepo,
	}
}

// CreateClass creates a new class
func (s *ClassServiceImpl) CreateClass(name string, startDate, endDate time.Time, capacity int) (models.Class, error) {
	class := models.NewClass(name, startDate, endDate, capacity)
	err := s.classRepo.Create(class)
	return class, err
}
