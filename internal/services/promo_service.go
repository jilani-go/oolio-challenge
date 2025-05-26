package services

import (
	"errors"
	"github.com/jilani-go/glofox/internal/repository"
	"log"
	"sync"
)

// Errors for PromoService
var (
	ErrInvalidPromoCode = errors.New("invalid promo code")
)

// PromoService defines the interface for promo code business logic
type PromoService interface {
	// ValidatePromoCode checks if a promo code is valid in a specific file
	ValidatePromoCode(code string) (bool, error)
}

// PromoServiceImpl implements PromoService
type PromoServiceImpl struct {
	promoRepo repository.PromoRepository
}

// NewPromoService creates a new promo service
func NewPromoService(promoRepo repository.PromoRepository) PromoService {
	return &PromoServiceImpl{
		promoRepo: promoRepo,
	}
}

// ValidatePromoCode checks if a promo code is valid in a specific file
func (s *PromoServiceImpl) ValidatePromoCode(code string) (bool, error) {
	if code == "" {
		return true, nil
	}
	if len(code) < 8 || len(code) > 10 {
		return false, ErrInvalidPromoCode
	}
	type result struct {
		exists bool
		err    error
	}
	results := make([]result, 3)
	var wg sync.WaitGroup

	// Launch goroutines to check files concurrently
	for i, fileNum := range []int{1, 2, 3} {
		wg.Add(1)
		go func(i, fileNum int) {
			defer wg.Done()
			exists, err := s.promoRepo.ExistsInFile(code, fileNum)
			results[i] = result{exists, err}
		}(i, fileNum)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	promoExistCount := 0
	for _, r := range results {
		if r.err != nil {
			return false, r.err //return if we get a single error
		}
		if r.exists {
			promoExistCount++
		}
	}
	log.Println("Promo code exists in ", promoExistCount, " files.")
	if promoExistCount >= 2 {
		return true, nil
	}

	return false, nil
}
