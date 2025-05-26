package models

import (
	"time"
)

// Product represents a food product in the system
type Product struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	Name        string    `json:"name" bson:"name"`
	Price       float64   `json:"price" bson:"price"`
	Category    string    `json:"category" bson:"category"`
	Description string    `json:"description,omitempty" bson:"description,omitempty"`
	ImageURL    string    `json:"imageUrl,omitempty" bson:"imageUrl,omitempty"`
	CreatedAt   time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt" bson:"updatedAt"`
}

