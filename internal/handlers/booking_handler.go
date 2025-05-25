package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/jilani-go/glofox/internal/constants"
	"github.com/jilani-go/glofox/internal/models"
	"github.com/jilani-go/glofox/internal/services"
	"github.com/jilani-go/glofox/internal/utils"
)

// BookingHandler handles booking-related HTTP requests
type BookingHandler struct {
	bookingService services.BookingService
}

// NewBookingHandler creates a new BookingHandler
func NewBookingHandler(bookingService services.BookingService) *BookingHandler {
	return &BookingHandler{
		bookingService: bookingService,
	}
}

// CreateBooking handles the creation of a new booking
func (h *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	var req BookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the request
	if err := utils.ValidateStruct(req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(utils.FormatValidationErrors(err))
		return
	}

	date, err := time.Parse(reqDateFormat, req.Date)
	if err != nil {
		http.Error(w, constants.ErrInvalidDateFormat.Error(), http.StatusBadRequest)
		return
	}
	booking, err := h.bookingService.CreateBooking(models.ClassID(req.ClassID), req.MemberName, date)
	if err != nil {
		if errors.Is(err, constants.ErrClassNotFound) || errors.Is(err, constants.ErrDateOutOfRange) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Booking created successfully",
		"id":      string(booking.ID),
	})
}
