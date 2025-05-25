package services

import (
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/jilani-go/glofox/internal/mocks"
	"github.com/jilani-go/glofox/internal/models"
	"github.com/stretchr/testify/assert"
)

// TestCreateClass tests the CreateClass method
func TestCreateClass(t *testing.T) {
	// Define test cases
	tests := []struct {
		name           string
		className      string
		startDate      time.Time
		endDate        time.Time
		capacity       int
		createErr      error
		expectedError  error
		expectedResult models.Class
	}{
		{
			name:      "Success - Create class",
			className: "Yoga",
			startDate: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			endDate:   time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC),
			capacity:  20,
			createErr: nil,
			expectedResult: models.Class{
				ID:        "00000000-0000-0000-0000-000000000001",
				Name:      "Yoga",
				StartDate: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC),
				Capacity:  20,
			},
		},
		{
			name:          "Failure - Error creating class",
			className:     "Yoga",
			startDate:     time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			endDate:       time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC),
			capacity:      20,
			createErr:     errors.New("database error"),
			expectedError: errors.New("database error"),
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup controller and mocks
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClassRepo := mocks.NewMockClassRepository(ctrl)
			mockBookingRepo := mocks.NewMockBookingRepository(ctrl)

			// Set expectations
			mockClassRepo.EXPECT().Create(gomock.Any()).DoAndReturn(
				func(class models.Class) error {
					// Verify class properties
					assert.Equal(t, tt.className, class.Name)
					assert.Equal(t, tt.startDate, class.StartDate)
					assert.Equal(t, tt.endDate, class.EndDate)
					assert.Equal(t, tt.capacity, class.Capacity)
					return tt.createErr
				})

			// Create service with mocks
			service := NewClassService(mockClassRepo, mockBookingRepo)

			// Call the method under test
			class, err := service.CreateClass(tt.className, tt.startDate, tt.endDate, tt.capacity)

			// Assert results
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult.Name, class.Name)
				assert.Equal(t, tt.expectedResult.StartDate, class.StartDate)
				assert.Equal(t, tt.expectedResult.EndDate, class.EndDate)
				assert.Equal(t, tt.expectedResult.Capacity, class.Capacity)
			}
		})
	}
}
