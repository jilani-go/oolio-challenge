package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jilani-go/glofox/internal/constants"
	"github.com/jilani-go/glofox/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockBookingService mocks the BookingService interface for testing
type MockBookingService struct {
	mock.Mock
}

func (m *MockBookingService) CreateBooking(classID models.ClassID, memberName string, bookingDate time.Time) (models.Booking, error) {
	args := m.Called(classID, memberName, bookingDate)
	return args.Get(0).(models.Booking), args.Error(1)
}

// TestCreateBooking tests the CreateBooking handler
func TestCreateBooking(t *testing.T) {
	// Define test cases
	tests := []struct {
		name            string
		requestBody     interface{}
		setupMock       func(*MockBookingService)
		expectedStatus  int
		expectedMessage string
	}{
		{
			name: "Success - Valid booking creation",
			requestBody: BookingRequest{
				ClassID:    "class-123",
				MemberName: "John Doe",
				Date:       "2023-01-15",
			},
			setupMock: func(mockService *MockBookingService) {
				expectedDate, _ := time.Parse("2006-01-02", "2023-01-15")
				mockService.On("CreateBooking",
					models.ClassID("class-123"),
					"John Doe",
					expectedDate,
				).Return(models.Booking{
					ID:          "booking-123",
					ClassID:     "class-123",
					MemberName:  "John Doe",
					BookingDate: expectedDate,
				}, nil)
			},
			expectedStatus:  http.StatusCreated,
			expectedMessage: "Booking created successfully",
		},
		{
			name: "Failure - Invalid request (missing member name)",
			requestBody: BookingRequest{
				ClassID: "class-123",
				// MemberName is missing
				Date: "2023-01-15",
			},
			setupMock:       func(mockService *MockBookingService) {},
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "MemberName is required",
		},
		{
			name: "Failure - Class not found",
			requestBody: BookingRequest{
				ClassID:    "non-existent",
				MemberName: "John Doe",
				Date:       "2023-01-15",
			},
			setupMock: func(mockService *MockBookingService) {
				expectedDate, _ := time.Parse("2006-01-02", "2023-01-15")
				mockService.On("CreateBooking",
					models.ClassID("non-existent"),
					"John Doe",
					expectedDate,
				).Return(models.Booking{}, constants.ErrClassNotFound)
			},
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: constants.ErrClassNotFound.Error(),
		},
		// Test case for class full has been removed since we no longer return an error for full classes
		{
			name: "Failure - Date out of range",
			requestBody: BookingRequest{
				ClassID:    "class-123",
				MemberName: "John Doe",
				Date:       "2023-01-15",
			},
			setupMock: func(mockService *MockBookingService) {
				expectedDate, _ := time.Parse("2006-01-02", "2023-01-15")
				mockService.On("CreateBooking",
					models.ClassID("class-123"),
					"John Doe",
					expectedDate,
				).Return(models.Booking{}, constants.ErrDateOutOfRange)
			},
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: constants.ErrDateOutOfRange.Error(),
		},
		{
			name: "Failure - Service error",
			requestBody: BookingRequest{
				ClassID:    "class-123",
				MemberName: "John Doe",
				Date:       "2023-01-15",
			},
			setupMock: func(mockService *MockBookingService) {
				expectedDate, _ := time.Parse("2006-01-02", "2023-01-15")
				mockService.On("CreateBooking",
					models.ClassID("class-123"),
					"John Doe",
					expectedDate,
				).Return(models.Booking{}, errors.New("service error"))
			},
			expectedStatus:  http.StatusInternalServerError,
			expectedMessage: "service error",
		},
		{
			name:            "Failure - Invalid JSON",
			requestBody:     "invalid json",
			setupMock:       func(mockService *MockBookingService) {},
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "Invalid request body",
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock service
			mockService := new(MockBookingService)

			// Setup mock expectations
			tt.setupMock(mockService)

			// Create handler
			handler := NewBookingHandler(mockService)

			// Create request
			var reqBody []byte
			var err error

			if str, ok := tt.requestBody.(string); ok {
				reqBody = []byte(str)
			} else {
				reqBody, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			req, err := http.NewRequest("POST", "/api/bookings", bytes.NewBuffer(reqBody))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call handler
			handler.CreateBooking(rr, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// Check response body contains expected message
			assert.Contains(t, rr.Body.String(), tt.expectedMessage)

			// Verify mock expectations
			mockService.AssertExpectations(t)
		})
	}
}
