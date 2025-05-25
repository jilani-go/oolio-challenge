package models

import (
	"time"

	"github.com/google/uuid"
)

type BookingID string
type Booking struct {
	ID          BookingID `json:"id"`
	ClassID     ClassID   `json:"class_id"`
	MemberName  string    `json:"member_name"`
	BookingDate time.Time `json:"booking_date"`
}

// NewBooking creates a new booking with a unique ID
func NewBooking(classID ClassID, memberName string, bookingDate time.Time) Booking {
	return Booking{
		ID:          BookingID(uuid.New().String()),
		ClassID:     classID,
		MemberName:  memberName,
		BookingDate: bookingDate,
	}
}
