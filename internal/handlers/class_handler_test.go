package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jilani-go/glofox/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockClassService mocks the ClassService interface for testing
type MockClassService struct {
	mock.Mock
}

func (m *MockClassService) CreateClass(name string, startDate, endDate time.Time, capacity int) (models.Class, error) {
	args := m.Called(name, startDate, endDate, capacity)
	return args.Get(0).(models.Class), args.Error(1)
}

// TestCreateClass tests the CreateClass handler
func TestCreateClass(t *testing.T) {
	// Define test cases
	tests := []struct {
		name            string
		requestBody     interface{}
		setupMock       func(*MockClassService)
		expectedStatus  int
		expectedMessage string
	}{
		{
			name: "Success - Valid class creation",
			requestBody: ClassRequest{
				Name:      "Yoga",
				StartDate: "2023-01-01",
				EndDate:   "2023-01-31",
				Capacity:  20,
			},
			setupMock: func(mockService *MockClassService) {
				startDate, _ := time.Parse("2006-01-02", "2023-01-01")
				endDate, _ := time.Parse("2006-01-02", "2023-01-31")
				mockService.On("CreateClass",
					"Yoga",
					startDate,
					endDate,
					20,
				).Return(models.Class{
					ID:        "class-123",
					Name:      "Yoga",
					StartDate: startDate,
					EndDate:   endDate,
					Capacity:  20,
				}, nil)
			},
			expectedStatus:  http.StatusCreated,
			expectedMessage: "Class created successfully",
		},
		{
			name: "Failure - Invalid request (missing name)",
			requestBody: ClassRequest{
				// Name field is missing
				StartDate: "2023-01-01",
				EndDate:   "2023-01-31",
				Capacity:  20,
			},
			setupMock:       func(mockService *MockClassService) {},
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "Name is required",
		},
		{
			name: "Failure - Invalid request (end date before start date)",
			requestBody: ClassRequest{
				Name:      "Yoga",
				StartDate: "2023-01-31",
				EndDate:   "2023-01-01", // Before start date
				Capacity:  20,
			},
			setupMock:       func(mockService *MockClassService) {},
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "EndDate must be after StartDate",
		},
		{
			name: "Failure - Invalid request (zero capacity)",
			requestBody: ClassRequest{
				Name:      "Yoga",
				StartDate: "2023-01-01",
				EndDate:   "2023-01-31",
				Capacity:  0, // Invalid capacity
			},
			setupMock:       func(mockService *MockClassService) {},
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "Capacity is required",
		},
		{
			name: "Failure - Service error",
			requestBody: ClassRequest{
				Name:      "Yoga",
				StartDate: "2023-01-01",
				EndDate:   "2023-01-31",
				Capacity:  20,
			},
			setupMock: func(mockService *MockClassService) {
				startDate, _ := time.Parse("2006-01-02", "2023-01-01")
				endDate, _ := time.Parse("2006-01-02", "2023-01-31")
				mockService.On("CreateClass",
					"Yoga",
					startDate,
					endDate,
					20,
				).Return(models.Class{}, errors.New("service error"))
			},
			expectedStatus:  http.StatusInternalServerError,
			expectedMessage: "service error",
		},
		{
			name:            "Failure - Invalid JSON",
			requestBody:     "invalid json",
			setupMock:       func(mockService *MockClassService) {},
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "Invalid request body",
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock service
			mockService := new(MockClassService)

			// Setup mock expectations
			tt.setupMock(mockService)

			// Create handler
			handler := NewClassHandler(mockService)

			// Create request
			var reqBody []byte
			var err error

			if str, ok := tt.requestBody.(string); ok {
				reqBody = []byte(str)
			} else {
				reqBody, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			req, err := http.NewRequest("POST", "/api/classes", bytes.NewBuffer(reqBody))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call handler
			handler.CreateClass(rr, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// Check response body contains expected message
			assert.Contains(t, rr.Body.String(), tt.expectedMessage)

			// Verify mock expectations
			mockService.AssertExpectations(t)
		})
	}
}
