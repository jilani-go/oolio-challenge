package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jilani-go/glofox/internal/handlers"
)

// SetupRoutes initializes the API routes
func SetupRoutes(classHandler *handlers.ClassHandler, bookingHandler *handlers.BookingHandler) http.Handler {
	// Create router
	router := mux.NewRouter()

	// Class routes
	router.HandleFunc("/api/classes", classHandler.CreateClass).Methods("POST")

	// Booking routes
	router.HandleFunc("/api/bookings", bookingHandler.CreateBooking).Methods("POST")

	return router
}
