package repository

import (
	"github.com/jilani-go/glofox/internal/models"
	"time"
)

//go:generate mockgen -destination=../mocks/mock_class_repository.go -package=mocks github.com/jilani-go/glofox/internal/repository ClassRepository
//go:generate mockgen -destination=../mocks/mock_booking_repository.go -package=mocks github.com/jilani-go/glofox/internal/repository BookingRepository

// ClassRepository defines operations for classes
type ClassRepository interface {
	Create(class models.Class) error
	GetByID(id models.ClassID) (models.Class, bool)
}

// BookingRepository defines operations for bookings
type BookingRepository interface {
	Create(booking models.Booking) error
	CountBookingsForClass(classID models.ClassID, date time.Time) (int, error)
}
