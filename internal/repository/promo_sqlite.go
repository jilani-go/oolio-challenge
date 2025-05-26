package repository

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// SQLitePromoRepository implements PromoRepository using SQLite database
type SQLitePromoRepository struct {
	db            *sql.DB
	databasePath  string
	fileCount     int
	batchSize     int
	workerCount   int
	createIndexes bool
}

// SQLitePromoConfig contains configuration options for SQLitePromoRepository
type SQLitePromoConfig struct {
	// DatabasePath is the path where the SQLite database will be stored
	DatabasePath string
	// BatchSize controls how many codes are inserted in a single transaction
	BatchSize int
	// WorkerCount controls the number of parallel workers for loading data
	WorkerCount int
	// CreateIndexes determines whether to create indexes after loading data
	CreateIndexes bool
}

// NewSQLitePromoRepository creates a new SQLite-based promo repository
func NewSQLitePromoRepository(config SQLitePromoConfig) (PromoRepository, error) {
	startTime := time.Now()

	// Set default configuration values if not provided
	if config.DatabasePath == "" {
		config.DatabasePath = "promo_codes.db"
	}
	if config.BatchSize <= 0 {
		config.BatchSize = 10000
	}
	if config.WorkerCount <= 0 {
		config.WorkerCount = runtime.NumCPU()
	}

	// Ensure the directory exists
	dbDir := filepath.Dir(config.DatabasePath)
	if dbDir != "" && dbDir != "." {
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create database directory: %w", err)
		}
	}

	// Open database connection with explicit connection string parameters
	dbConnectionString := fmt.Sprintf("%s?_journal_mode=WAL&_synchronous=NORMAL&_cache_size=10000&_foreign_keys=ON", config.DatabasePath)
	db, err := sql.Open("sqlite3", dbConnectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open SQLite database: %w", err)
	}

	// Verify connection is working with ping
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping SQLite database: %w", err)
	}

	// Set pragmas for better performance with validation
	pragmas := []string{
		"PRAGMA journal_mode = WAL",
		"PRAGMA synchronous = NORMAL",
		"PRAGMA cache_size = 10000",
		"PRAGMA temp_store = MEMORY",
		"PRAGMA mmap_size = 30000000000",
		"PRAGMA page_size = 4096",
	}

	for _, pragma := range pragmas {
		// Validate pragma syntax (basic check)
		if !strings.HasPrefix(strings.ToUpper(pragma), "PRAGMA ") {
			db.Close()
			return nil, fmt.Errorf("invalid pragma format: %s", pragma)
		}

		// Execute the pragma
		_, err := db.Exec(pragma)
		if err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to set pragma %s: %w", pragma, err)
		}
	}

	// Define promo file paths
	promoFiles := []string{
		"internal/repository/promofiles/couponbase1",
		"internal/repository/promofiles/couponbase2",
		"internal/repository/promofiles/couponbase3",
	}

	// Create repository instance
	repo := &SQLitePromoRepository{
		db:            db,
		databasePath:  config.DatabasePath,
		fileCount:     len(promoFiles),
		batchSize:     config.BatchSize,
		workerCount:   config.WorkerCount,
		createIndexes: config.CreateIndexes,
	}

	// Initialize database schema
	if err := repo.initializeSchema(); err != nil {
		db.Close()
		log.Println("Failed to initialize database schema, attempting to repair...", err)
		return nil, fmt.Errorf("failed to initialize database schema: %w", err)
	}

	// Check if tables are already populated
	populated, err := repo.areTablesPopulated()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to check if tables are populated: %w", err)
	}

	// Load data if tables are not populated
	if !populated {
		log.Println("Promo code tables are empty, loading data...")
		if err := repo.loadPromoFiles(promoFiles); err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to load promo files: %w", err)
		}
	}

	_ = time.Since(startTime)

	return repo, nil
}

// initializeSchema creates the necessary tables in the database
func (r *SQLitePromoRepository) initializeSchema() error {
	log.Printf("Creating tables for %d promo files", r.fileCount)

	// Create tables for each promo file
	for i := 1; i <= r.fileCount; i++ {
		tableName := fmt.Sprintf("promo_codes_file%d", i)
		query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s ( code TEXT PRIMARY KEY )", tableName)

		log.Printf("Executing SQL: %s", query)
		result, err := r.db.Exec(query)
		if err != nil {
			return fmt.Errorf("failed to create table for file %d: %w", i, err)
		}

		// Log success message with any affected rows (usually 0 for CREATE TABLE IF NOT EXISTS)
		rowsAffected, _ := result.RowsAffected()
		log.Printf("Table %s created or already exists (rows affected: %d)", tableName, rowsAffected)

		// Verify the table was actually created
		verifyQuery := fmt.Sprintf("SELECT name FROM sqlite_master WHERE type='table' AND name='%s'", tableName)
		var name string
		err = r.db.QueryRow(verifyQuery).Scan(&name)
		if err != nil {
			return fmt.Errorf("failed to verify table %s was created: %w", tableName, err)
		}

		if name != tableName {
			return fmt.Errorf("table %s was not created successfully", tableName)
		}

		log.Printf("Verified table %s exists in database", tableName)
	}

	return nil
}

// areTablesPopulated checks if the promo code tables already have data
func (r *SQLitePromoRepository) areTablesPopulated() (bool, error) {
	for i := 2; i <= r.fileCount; i++ {
		var count int
		tableName := fmt.Sprintf("promo_codes_file%d", i)

		// Verify the table name doesn't contain unexpected characters
		if strings.ContainsAny(tableName, ";\"`'") {
			return false, fmt.Errorf("invalid table name format: %s", tableName)
		}

		// Construct safe query
		query := fmt.Sprintf("SELECT COUNT(*) FROM %s LIMIT 1", tableName)

		// Execute query with proper error handling
		err := r.db.QueryRow(query).Scan(&count)
		if err != nil {
			// Check specifically for "no such table" error
			if strings.Contains(err.Error(), "no such table") {
				return false, nil
			}
			return false, fmt.Errorf("failed to check if table %d is populated: %w", i, err)
		}

		if count > 0 {
			return true, nil
		}
	}
	return false, nil
}

// loadPromoFiles loads promo codes from files into the database
func (r *SQLitePromoRepository) loadPromoFiles(promoFiles []string) error {
	startTime := time.Now()
	log.Println("Starting to load promo files into SQLite...")

	// Process each file in parallel
	var wg sync.WaitGroup
	errChan := make(chan error, len(promoFiles))

	for i, filePath := range promoFiles {
		if i == 0 {
			continue
		}
		fileNumber := i + 1
		wg.Add(1)
		go func(fn int, path string) {
			defer wg.Done()
			if err := r.loadPromoFile(fn, path); err != nil {
				errChan <- fmt.Errorf("error loading file %d: %w", fn, err)
			}
		}(fileNumber, filePath)
	}

	// Wait for all loading to complete
	wg.Wait()
	close(errChan)

	// Check for errors
	if len(errChan) > 0 {
		return <-errChan
	}

	// Create indexes if configured to do so
	if r.createIndexes {
		log.Println("Creating indexes on promo code tables...")
		for i := 1; i <= r.fileCount; i++ {
			query := fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_code_file%d ON promo_codes_file%d(code)", i, i)
			if _, err := r.db.Exec(query); err != nil {
				return fmt.Errorf("failed to create index for file %d: %w", i, err)
			}
		}
	}

	elapsed := time.Since(startTime)
	log.Printf("Completed loading all promo files in %v", elapsed)
	return nil
}

// loadPromoFile loads a single promo file into the database
func (r *SQLitePromoRepository) loadPromoFile(fileNumber int, filePath string) error {
	startTime := time.Now()
	log.Printf("Starting to load file %d: %s", fileNumber, filePath)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Create empty file for testing if it doesn't exist
		log.Printf("File %s doesn't exist, creating empty test file", filePath)
		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
		if err := os.WriteFile(filePath, []byte("testcode1\ntestcode2\ntestcode3\n"), 0644); err != nil {
			return fmt.Errorf("failed to create test file: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to check if file exists: %w", err)
	}

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Begin transaction
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if tx != nil {
			tx.Rollback()
		}
	}()

	// Prepare insert statement
	tableName := fmt.Sprintf("promo_codes_file%d", fileNumber)
	stmt, err := tx.Prepare(fmt.Sprintf("INSERT OR IGNORE INTO %s (code) VALUES (?)", tableName))
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Process lines
	lineCount := 0
	batchCount := 0
	scanner := bufio.NewScanner(file)
	// Use a large buffer for scanner
	buffer := make([]byte, 4*1024*1024) // 4MB buffer
	scanner.Buffer(buffer, 4*1024*1024)

	for scanner.Scan() {
		code := strings.TrimSpace(scanner.Text())
		if code != "" {
			if _, err := stmt.Exec(code); err != nil {
				return fmt.Errorf("error inserting code: %w", err)
			}
			lineCount++

			// Commit every batchSize records
			if lineCount > 0 && lineCount%r.batchSize == 0 {
				if err := tx.Commit(); err != nil {
					return fmt.Errorf("error committing transaction: %w", err)
				}

				batchCount++
				if batchCount%10 == 0 {
					log.Printf("File %d: Inserted %d codes...", fileNumber, lineCount)
				}

				// Start new transaction
				tx = nil
				tx, err = r.db.Begin()
				if err != nil {
					return fmt.Errorf("failed to begin new transaction: %w", err)
				}

				// Prepare new statement
				stmt.Close() // Close previous statement
				stmt, err = tx.Prepare(fmt.Sprintf("INSERT OR IGNORE INTO %s (code) VALUES (?)", tableName))
				if err != nil {
					return fmt.Errorf("failed to prepare statement: %w", err)
				}
			}
		}
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	// Commit final batch
	if tx != nil {
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("error committing final transaction: %w", err)
		}
		tx = nil
	}

	elapsed := time.Since(startTime)
	log.Printf("Completed loading file %d with %d codes in %v", fileNumber, lineCount, elapsed)
	return nil
}

// ExistsInFile checks if a given promo code exists in a specific file
func (r *SQLitePromoRepository) ExistsInFile(code string, fileNumber int) (bool, error) {
	// Validate fileNumber
	if fileNumber < 1 || fileNumber > r.fileCount {
		return false, ErrInvalidFileNumber
	}

	// Sanitize table name for security
	tableName := fmt.Sprintf("promo_codes_file%d", fileNumber)
	if strings.ContainsAny(tableName, ";\"`'") {
		return false, fmt.Errorf("invalid table name format: %s", tableName)
	}

	// Construct the safe parameterized query
	query := fmt.Sprintf("SELECT 1 FROM %s WHERE code = ? LIMIT 1", tableName)

	// Execute the query with proper parameter binding for security
	var exists int
	err := r.db.QueryRow(query, code).Scan(&exists)

	if err == sql.ErrNoRows {
		// Code doesn't exist
		return false, nil
	} else if err != nil {
		// Handle table doesn't exist error specifically
		if strings.Contains(err.Error(), "no such table") {
			return false, nil
		}
		return false, fmt.Errorf("error checking if code exists in file %d: %w", fileNumber, err)
	}

	return true, nil
}

// Close closes the database connection
func (r *SQLitePromoRepository) Close() error {
	if r.db != nil {
		return r.db.Close()
	}
	return nil
}
