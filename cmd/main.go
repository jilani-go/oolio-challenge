package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/jilani-go/glofox/internal/api"
	"github.com/jilani-go/glofox/internal/config"
	"github.com/jilani-go/glofox/internal/handlers"
	"github.com/jilani-go/glofox/internal/repository"
	"github.com/jilani-go/glofox/internal/services"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Create repositories
	productRepo := repository.NewInMemoryProductRepository()
	orderRepo := repository.NewInMemoryOrderRepository(productRepo)

	// Create SQLite promo repository with optimized configuration
	promoRepo, err := repository.NewSQLitePromoRepository(repository.SQLitePromoConfig{
		DatabasePath:  "data/promo_codes.db",
		BatchSize:     50000, // Insert 50k records per transaction
		WorkerCount:   runtime.NumCPU(),
		CreateIndexes: true, // Create indexes for faster lookups
	})
	if err != nil {
		log.Fatalf("Failed to initialize promo repository: %v", err)
	}

	// Create services
	productService := services.NewProductService(productRepo)
	promoService := services.NewPromoService(promoRepo)
	orderService := services.NewOrderService(orderRepo, productRepo)

	// Create handlers
	productHandler := handlers.NewProductHandler(productService)
	orderHandler := handlers.NewOrderHandler(orderService, productService, promoService)

	// Setup routes
	router := api.SetupRoutes(productHandler, orderHandler)

	// Override port from environment if provided
	if envPort := os.Getenv("PORT"); envPort != "" {
		cfg.Server.Port = envPort
	}

	// Configure HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Channel to listen for errors coming from the listener
	serverErrors := make(chan error, 1)

	// Start the server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.Server.Port)
		serverErrors <- server.ListenAndServe()
	}()

	// Channel to listen for interrupt/terminate signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Block until an os.Signal or an error is received
	select {
	case err := <-serverErrors:
		log.Fatalf("Error starting server: %v", err)

	case sig := <-shutdown:
		log.Printf("Shutdown signal received: %v", sig)

		// Create a deadline context for graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Gracefully shutdown connections
		log.Println("Shutting down server...")

		// Shut down database connections
		if err := promoRepo.Close(); err != nil {
			log.Printf("Error closing SQLite connection: %v", err)
		}

		// Shut down HTTP server
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Error during server shutdown: %v", err)
			server.Close()
		}

		// Verify if the server shutdown gracefully
		<-ctx.Done()
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Shutdown deadline exceeded, forcing exit")
		} else {
			log.Println("Server gracefully stopped")
		}
	}
}
