package repository

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

// Constants for optimized file reading
const (
	// readerBufferSize defines the buffer size for file reading
	readerBufferSize = 1 * 1024 * 1024 // 1MB buffer

	// workerChannelBuffer defines the buffer size for the worker channel
	workerChannelBuffer = 200_000 // Buffer for 200,000 lines

	// initialMapCapacity provides an initial capacity for the promo code maps
	initialMapCapacity = 10_000_000 // Expect ~10 million unique codes per file

	// numWorkersMultiplier determines how many worker goroutines to create per CPU core
	numWorkersMultiplier = 2
)

// Errors for PromoRepository
var (
	ErrInvalidFileNumber = errors.New("invalid file number")
	ErrLoadingPromoFile  = errors.New("error loading promo file")
)

// PromoRepository defines the interface for promo code operations
type PromoRepository interface {
	// ExistsInFile checks if a given promo code exists in a specific file
	ExistsInFile(code string, fileNumber int) (bool, error)

	// Close cleans up any resources used by the repository
	// Implementation is optional for repositories that don't require cleanup
	Close() error
}

// InMemoryPromoRepository implements PromoRepository using in-memory maps
type InMemoryPromoRepository struct {
	// Each key represents a file number, each value is a map of promo codes in that file
	filePromoCodes map[int]map[string]struct{}
}

// Global mutex for synchronizing access to shared data
var mu sync.Mutex

// NewInMemoryPromoRepository creates a new in-memory promo repository
func NewInMemoryPromoRepository() (PromoRepository, error) {
	startTime := time.Now()
	log.Println("Starting parallel loading of promo code files...")

	// Define promo file paths
	promoFiles := []string{
		"internal/repository/promofiles/couponbase1",
		"internal/repository/promofiles/couponbase2",
		"internal/repository/promofiles/couponbase3",
	}

	// Create a channel for results
	type result struct {
		fileNumber int
		promoCodes map[string]struct{}
		err        error
	}
	resultChan := make(chan result, len(promoFiles))

	// Process files in parallel
	var wg sync.WaitGroup
	for i, filePath := range promoFiles {
		fileNumber := i + 1
		wg.Add(1)

		// Launch a goroutine to process each file
		go func(fn int, path string) {
			defer wg.Done()

			fileStartTime := time.Now()
			log.Printf("Starting to load file %d: %s", fn, path)

			// Process the file
			promoCodes, err := readPromoCodesOptimized(path)

			// Send result back through channel
			if err != nil {
				log.Printf("Error loading file %d: %v", fn, err)
				resultChan <- result{fn, nil, fmt.Errorf("failed to load file %d: %w", fn, err)}
			} else {
				elapsed := time.Since(fileStartTime)
				log.Printf("Completed loading file %d (%s) with %d codes in %v",
					fn, path, len(promoCodes), elapsed)
				resultChan <- result{fn, promoCodes, nil}
			}
		}(fileNumber, filePath)
	}

	// Close result channel when all files are processed
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results and build the final map
	filePromoCodes := make(map[int]map[string]struct{})
	for r := range resultChan {
		if r.err != nil {
			return nil, r.err
		}
		filePromoCodes[r.fileNumber] = r.promoCodes
	}

	totalElapsed := time.Since(startTime)
	log.Printf("Successfully loaded all promo files in %v", totalElapsed)

	return &InMemoryPromoRepository{
		filePromoCodes: filePromoCodes,
	}, nil
}

// readPromoCodesFromFile provides backward compatibility
func readPromoCodesFromFile(filePath string) (map[string]struct{}, error) {
	return readPromoCodesOptimized(filePath)
}

// readPromoCodesOptimized reads promo codes from a file using concurrent workers
func readPromoCodesOptimized(filePath string) (map[string]struct{}, error) {
	startTime := time.Now()

	// --- Step 1: Open the file ---
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// --- Step 2: Initialize buffered reader ---
	reader := bufio.NewReaderSize(file, readerBufferSize)

	// --- Step 3: Prepare for concurrency ---
	// Channel for distributing lines to workers
	linesChan := make(chan string, workerChannelBuffer)

	// Pre-allocate map with reasonable capacity
	globalMap := make(map[string]struct{}, initialMapCapacity)

	// Mutex for protecting access to the global map
	var mapMutex sync.Mutex

	// WaitGroup for coordinating goroutines
	var wg sync.WaitGroup

	// Error channel for propagating errors from goroutines
	errChan := make(chan error, 1)

	// --- Step 4: Determine the number of workers ---
	numWorkers := runtime.NumCPU() * numWorkersMultiplier
	if numWorkers < 2 {
		numWorkers = 2 // Ensure at least 2 workers
	}

	// --- Step 5: Reader Goroutine (Fan-out) ---
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(linesChan)

		lineCount := 0
		lastReport := time.Now()
		reportInterval := 5 * time.Second

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break // End of file
				}
				// Report error and exit
				select {
				case errChan <- fmt.Errorf("error reading file: %w", err):
				default:
					// Channel already has an error
				}
				return
			}

			// Remove newline and trim whitespace
			trimmed := strings.TrimSpace(line)
			if trimmed != "" {
				// Send the line to workers
				linesChan <- trimmed
			}

			// Track progress
			lineCount++
			if lineCount%1000000 == 0 { // Every million lines
				now := time.Now()
				if now.Sub(lastReport) >= reportInterval {
					log.Printf("File %s: Read %d million lines...", filePath, lineCount/1000000)
					lastReport = now
				}
			}
		}
	}()

	// --- Step 6: Worker Goroutines (Fan-in) ---
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			// Each worker maintains its own local map
			localMap := make(map[string]struct{})

			// Process lines from channel
			for line := range linesChan {
				localMap[line] = struct{}{}
			}

			// Merge local map into the global map
			mapMutex.Lock()
			for k := range localMap {
				globalMap[k] = struct{}{}
			}
			mapMutex.Unlock()

			// Help with garbage collection
			localMap = nil
		}(i)
	}

	// --- Step 7: Wait for all goroutines to finish ---
	wg.Wait()
	close(errChan)

	// Check for errors
	select {
	case err := <-errChan:
		if err != nil {
			return nil, err
		}
	default:
		// No errors
	}

	elapsed := time.Since(startTime)
	log.Printf("Completed reading %s with %d unique codes in %v",
		filePath, len(globalMap), elapsed)

	return globalMap, nil
}

// ExistsInFile checks if a given promo code exists in a specific file
func (r *InMemoryPromoRepository) ExistsInFile(code string, fileNumber int) (bool, error) {
	// Validate fileNumber
	if fileNumber < 1 || fileNumber > len(r.filePromoCodes) {
		return false, ErrInvalidFileNumber
	}

	// Get promo codes for the specified file
	promoCodes, exists := r.filePromoCodes[fileNumber]
	if !exists {
		return false, nil
	}

	// Check if code exists in the file
	_, codeExists := promoCodes[code]
	return codeExists, nil
}

// Close implements the Close method for PromoRepository interface
// For in-memory repository, this is a no-op
func (r *InMemoryPromoRepository) Close() error {
	return nil
}
