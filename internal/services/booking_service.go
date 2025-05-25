package services

import (
	"log"
	"time"

	"github.com/jilani-go/glofox/internal/constants"
	"github.com/jilani-go/glofox/internal/models"
	"github.com/jilani-go/glofox/internal/repository"
	"github.com/jilani-go/glofox/internal/utils"
)

// BookingService defines operations for booking business logic
type BookingService interface {
	CreateBooking(classID models.ClassID, memberName string, bookingDate time.Time) (models.Booking, error)
}

// BookingServiceImpl implements BookingService
type BookingServiceImpl struct {
	bookingRepo  repository.BookingRepository
	classRepo    repository.ClassRepository
	classService ClassService
}

// NewBookingService creates a new booking service
func NewBookingService(
	bookingRepo repository.BookingRepository,
	classRepo repository.ClassRepository,
	classService ClassService,
) BookingService {
	return &BookingServiceImpl{
		bookingRepo:  bookingRepo,
		classRepo:    classRepo,
		classService: classService,
	}
}

// CreateBooking creates a new booking for a class
func (s *BookingServiceImpl) CreateBooking(classID models.ClassID, memberName string, bookingDate time.Time) (models.Booking, error) {
	// Check if class exists
	class, exists := s.classRepo.GetByID(classID)
	if !exists {
		return models.Booking{}, constants.ErrClassNotFound
	}

	date := utils.NormalizeToMidnightUTC(bookingDate)

	// Check if booking date is within class schedule
	if date.Before(utils.NormalizeToMidnightUTC(class.StartDate)) || date.After(utils.NormalizeToMidnightUTC(class.EndDate)) {
		return models.Booking{}, constants.ErrDateOutOfRange
	}

	// Check available capacity
	currentCapacity, err := s.bookingRepo.CountBookingsForClass(classID, date)
	if err != nil {
		return models.Booking{}, err
	}

	if currentCapacity >= class.Capacity {
		log.Printf("Class %s is at full capacity (%d/%d) for date %s", classID, currentCapacity, class.Capacity, date.Format("2006-01-02"))
	}

	// Create booking
	booking := models.NewBooking(classID, memberName, date)
	err = s.bookingRepo.Create(booking)
	if err != nil {
		return models.Booking{}, err
	}

	return booking, nil
}
