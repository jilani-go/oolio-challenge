package handlers

// ClassRequest represents the API request for creating a class
type ClassRequest struct {
	Name      string `json:"name" validate:"required"`
	StartDate string `json:"start_date" validate:"required"`
	EndDate   string `json:"end_date" validate:"required,gtefield=StartDate"`
	Capacity  int    `json:"capacity" validate:"required,min=1"`
}

// BookingRequest represents the API request for creating a booking
type BookingRequest struct {
	ClassID    string `json:"class_id" validate:"required"`
	MemberName string `json:"member_name" validate:"required"`
	Date       string `json:"date" validate:"required"`
}

const (
	reqDateFormat = "2006-01-02"
)
