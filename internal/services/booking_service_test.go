package services

import (
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/jilani-go/glofox/internal/constants"
	"github.com/jilani-go/glofox/internal/mocks"
	"github.com/jilani-go/glofox/internal/models"
	"github.com/jilani-go/glofox/internal/utils"
	"github.com/stretchr/testify/assert"
)

// compareBookingIgnoreID is a helper function to compare two bookings ignoring the ID field
func compareBookingIgnoreID(t *testing.T, expected, actual models.Booking) {
	t.Helper()
	assert.Equal(t, expected.ClassID, actual.ClassID)
	assert.Equal(t, expected.MemberName, actual.MemberName)
	assert.True(t, expected.BookingDate.Equal(actual.BookingDate))
	assert.NotEmpty(t, actual.ID) // Still check that ID is not empty
}

// compareClassIgnoreID is a helper function to compare two classes ignoring the ID field
func compareClassIgnoreID(t *testing.T, expected, actual models.Class) {
	t.Helper()
	assert.Equal(t, expected.Name, actual.Name)
	assert.True(t, expected.StartDate.Equal(actual.StartDate))
	assert.True(t, expected.EndDate.Equal(actual.EndDate))
	assert.Equal(t, expected.Capacity, actual.Capacity)
	assert.NotEmpty(t, actual.ID) // Still check that ID is not empty
}

// createTestBooking is a helper function to create a booking for tests
func createTestBooking(classID models.ClassID, memberName string, bookingDate time.Time) models.Booking {
	return models.Booking{
		ID:          "", // ID will be ignored in comparison
		ClassID:     classID,
		MemberName:  memberName,
		BookingDate: bookingDate,
	}
}

// TestCreateBooking tests the CreateBooking method
func TestCreateBooking(t *testing.T) {

	// Define test cases
	tests := []struct {
		name           string
		classID        models.ClassID
		memberName     string
		bookingDate    time.Time
		classExists    bool
		class          models.Class
		bookingCount   int
		countErr       error
		createErr      error
		expectedError  error
		expectedResult models.Booking
	}{
		{
			name:        "Success - Create booking",
			classID:     "class-123",
			memberName:  "John Doe",
			bookingDate: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			classExists: true,
			class: models.Class{
				ID:        "class-123",
				Name:      "Yoga",
				StartDate: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC),
				Capacity:  20,
			},
			bookingCount: 5,
			countErr:     nil,
			createErr:    nil,
			expectedResult: models.Booking{
				ClassID:     "class-123",
				MemberName:  "John Doe",
				BookingDate: time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name:          "Failure - Class not found",
			classID:       "non-existent",
			memberName:    "John Doe",
			bookingDate:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			classExists:   false,
			expectedError: constants.ErrClassNotFound,
		},
		{
			name:        "Failure - Date out of range (before start)",
			classID:     "class-123",
			memberName:  "John Doe",
			bookingDate: time.Date(2022, 12, 15, 10, 0, 0, 0, time.UTC),
			classExists: true,
			class: models.Class{
				ID:        "class-123",
				Name:      "Yoga",
				StartDate: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC),
				Capacity:  20,
			},
			expectedError: constants.ErrDateOutOfRange,
		},
		{
			name:        "Failure - Date out of range (after end)",
			classID:     "class-123",
			memberName:  "John Doe",
			bookingDate: time.Date(2023, 2, 15, 10, 0, 0, 0, time.UTC),
			classExists: true,
			class: models.Class{
				ID:        "class-123",
				Name:      "Yoga",
				StartDate: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC),
				Capacity:  20,
			},
			expectedError: constants.ErrDateOutOfRange,
		},
		{
			name:        "Success - Class full but booking still created",
			classID:     "class-123",
			memberName:  "John Doe",
			bookingDate: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			classExists: true,
			class: models.Class{
				ID:        "class-123",
				Name:      "Yoga",
				StartDate: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC),
				Capacity:  20,
			},
			bookingCount: 20, // Equal to capacity
			expectedResult: models.Booking{
				ClassID:     "class-123",
				MemberName:  "John Doe",
				BookingDate: time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name:        "Failure - Error counting bookings",
			classID:     "class-123",
			memberName:  "John Doe",
			bookingDate: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			classExists: true,
			class: models.Class{
				ID:        "class-123",
				Name:      "Yoga",
				StartDate: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC),
				Capacity:  20,
			},
			countErr:      errors.New("database error"),
			expectedError: errors.New("database error"),
		},
		{
			name:        "Failure - Error creating booking",
			classID:     "class-123",
			memberName:  "John Doe",
			bookingDate: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			classExists: true,
			class: models.Class{
				ID:        "class-123",
				Name:      "Yoga",
				StartDate: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC),
				Capacity:  20,
			},
			bookingCount:  5,
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
			mockClassService := mocks.NewMockClassService(ctrl)

			// Set expectations
			mockClassRepo.EXPECT().GetByID(tt.classID).Return(tt.class, tt.classExists)

			if tt.classExists && !tt.bookingDate.Before(tt.class.StartDate) && !tt.bookingDate.After(tt.class.EndDate) {
				expectedDate := utils.NormalizeToMidnightUTC(tt.bookingDate)
				mockBookingRepo.EXPECT().CountBookingsForClass(tt.classID, expectedDate).Return(tt.bookingCount, tt.countErr)

				if tt.countErr == nil {
					// Match any booking with the expected properties
					mockBookingRepo.EXPECT().Create(gomock.Any()).DoAndReturn(
						func(booking models.Booking) error {
							// Verify booking properties
							assert.Equal(t, tt.classID, booking.ClassID)
							assert.Equal(t, tt.memberName, booking.MemberName)
							assert.True(t, expectedDate.Equal(booking.BookingDate))
							return tt.createErr
						})
				}
			}

			// Create service with mocks
			service := NewBookingService(mockBookingRepo, mockClassRepo, mockClassService)

			// Call the method under test
			booking, err := service.CreateBooking(tt.classID, tt.memberName, tt.bookingDate)

			// Assert results
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				// Use the custom comparison function that ignores ID
				compareBookingIgnoreID(t, tt.expectedResult, booking)
			}
		})
	}
}
