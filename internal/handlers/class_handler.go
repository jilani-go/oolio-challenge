package handlers

import (
	"encoding/json"
	"github.com/jilani-go/glofox/internal/constants"
	"net/http"
	"time"

	"github.com/jilani-go/glofox/internal/services"
	"github.com/jilani-go/glofox/internal/utils"
)

// ClassHandler handles class-related HTTP requests
type ClassHandler struct {
	classService services.ClassService
}

// NewClassHandler creates a new ClassHandler
func NewClassHandler(classService services.ClassService) *ClassHandler {
	return &ClassHandler{
		classService: classService,
	}
}

// CreateClass handles the creation of a new class
func (h *ClassHandler) CreateClass(w http.ResponseWriter, r *http.Request) {
	var req ClassRequest
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

	startDate, err := time.Parse(reqDateFormat, req.StartDate)
	if err != nil {
		http.Error(w, constants.ErrInvalidDateFormat.Error(), http.StatusBadRequest)
		return
	}

	endDate, err := time.Parse(reqDateFormat, req.EndDate)
	if err != nil {
		http.Error(w, constants.ErrInvalidDateFormat.Error(), http.StatusBadRequest)
		return
	}

	if startDate.After(endDate) {
		http.Error(w, "EndDate must be after StartDate", http.StatusBadRequest)
		return
	}
	// Call service to create a class
	class, err := h.classService.CreateClass(req.Name, startDate, endDate, req.Capacity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Class created successfully",
		"id":      string(class.ID),
	})
}
