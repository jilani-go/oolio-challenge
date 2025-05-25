package repository

import (
	"sync"
	"time"

	"github.com/jilani-go/glofox/internal/models"
	"github.com/jilani-go/glofox/internal/utils"
)

// InMemoryBookingRepo manages the in-memory booking data
type InMemoryBookingRepo struct {
	bookings map[models.ClassID]map[time.Time][]models.Booking
	mu       sync.RWMutex
}

// NewInMemoryBookingRepo creates a new InMemoryBookingRepo
func NewInMemoryBookingRepo() *InMemoryBookingRepo {
	return &InMemoryBookingRepo{
		bookings: make(map[models.ClassID]map[time.Time][]models.Booking),
	}
}

// Create for creating a new booking
func (bookingRepo *InMemoryBookingRepo) Create(booking models.Booking) error {
	bookingRepo.mu.Lock()
	defer bookingRepo.mu.Unlock()

	// Normalize date to midnight
	date := utils.NormalizeToMidnightUTC(booking.BookingDate)

	if _, exists := bookingRepo.bookings[booking.ClassID]; !exists {
		bookingRepo.bookings[booking.ClassID] = make(map[time.Time][]models.Booking)
	}

	bookingRepo.bookings[booking.ClassID][date] = append(bookingRepo.bookings[booking.ClassID][date], booking)
	return nil
}

// CountBookingsForClass counts bookings for a class on a specific date
func (bookingRepo *InMemoryBookingRepo) CountBookingsForClass(classID models.ClassID, date time.Time) (int, error) {
	bookingRepo.mu.RLock()
	defer bookingRepo.mu.RUnlock()

	classBookings, classExists := bookingRepo.bookings[classID]
	if !classExists {
		return 0, nil
	}

	dateBookings, dateExists := classBookings[date]
	if !dateExists {
		return 0, nil
	}

	return len(dateBookings), nil
}
