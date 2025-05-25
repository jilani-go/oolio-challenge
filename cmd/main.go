package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jilani-go/glofox/internal/api"
	"github.com/jilani-go/glofox/internal/config"
	"github.com/jilani-go/glofox/internal/handlers"
	"github.com/jilani-go/glofox/internal/repository"
	"github.com/jilani-go/glofox/internal/services"
)

type appDependencies struct {
	cfg            *config.Config
	router         http.Handler
	classHandler   *handlers.ClassHandler
	bookingHandler *handlers.BookingHandler
}

// initDependencies initializes all dependencies
func initDependencies() appDependencies {
	// Load configuration
	cfg := config.Load()

	// Initialize repositories
	classRepo := repository.NewInMemoryClassRepo()
	bookingRepo := repository.NewInMemoryBookingRepo()

	// Initialize services
	classService := services.NewClassService(classRepo, bookingRepo)
	bookingService := services.NewBookingService(bookingRepo, classRepo, classService)

	// Initialize handlers with services
	classHandler := handlers.NewClassHandler(classService)
	bookingHandler := handlers.NewBookingHandler(bookingService)

	// Setup router
	router := api.SetupRoutes(classHandler, bookingHandler)

	return appDependencies{
		cfg:            cfg,
		router:         router,
		classHandler:   classHandler,
		bookingHandler: bookingHandler,
	}
}

// setupServer creates and configures the HTTP server
func setupServer(deps appDependencies) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%s", deps.cfg.Server.Port),
		Handler:      deps.router,
		ReadTimeout:  deps.cfg.Server.ReadTimeout,
		WriteTimeout: deps.cfg.Server.WriteTimeout,
		IdleTimeout:  deps.cfg.Server.IdleTimeout,
	}
}

// gracefulShutdown handles graceful server shutdown with a timeout
func gracefulShutdown(server *http.Server, timeout time.Duration) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline to wait for current operations to complete
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}

func main() {
	deps := initDependencies()

	server := setupServer(deps)

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s...\n", deps.cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Handle graceful shutdown
	gracefulShutdown(server, 15*time.Second)
}
