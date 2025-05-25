package models

import (
	"time"

	"github.com/google/uuid"
)

type ClassID string
type Class struct {
	ID        ClassID   `json:"id"`
	Name      string    `json:"name"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Capacity  int       `json:"capacity"`
}

// NewClass creates a new class with a unique ID
func NewClass(name string, startDate, endDate time.Time, capacity int) Class {
	return Class{
		ID:        ClassID(uuid.New().String()),
		Name:      name,
		StartDate: startDate,
		EndDate:   endDate,
		Capacity:  capacity,
	}
}
