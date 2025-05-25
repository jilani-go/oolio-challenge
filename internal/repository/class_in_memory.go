package repository

import (
	"github.com/jilani-go/glofox/internal/models"
	"sync"
)

// InMemoryClassRepo manages the in-memory class data
type InMemoryClassRepo struct {
	classes     map[models.ClassID]models.Class
	mu          sync.RWMutex
	bookingRepo BookingRepository
}

// NewInMemoryClassRepo creates a new InMemoryClassRepo
func NewInMemoryClassRepo() *InMemoryClassRepo {
	return &InMemoryClassRepo{
		classes: make(map[models.ClassID]models.Class),
	}
}

// Create for creating a new class
func (classRepo *InMemoryClassRepo) Create(class models.Class) error {
	classRepo.mu.Lock()
	defer classRepo.mu.Unlock()

	classRepo.classes[class.ID] = class
	return nil
}

// GetByID gets a class by ID
func (classRepo *InMemoryClassRepo) GetByID(id models.ClassID) (models.Class, bool) {
	classRepo.mu.RLock()
	defer classRepo.mu.RUnlock()

	class, exists := classRepo.classes[id]
	return class, exists
}
