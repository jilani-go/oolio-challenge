# Promo Code Validation Service

A high-performance REST API service for validating promotion codes with optimized storage and retrieval capabilities. This system validates promotion codes against multiple data sources and determines their validity based on business rules.

## Features

- **Fast Promo Code Validation**: Validates codes against multiple source files simultaneously
- **Multi-Source Validation**: Requires code presence in at least 2 of 3 data sources for validity
- **Multiple Repository Implementations**:
  - In-memory repository for high-speed operation
  - SQLite repository for persistent storage with optimized configuration
- **Concurrent Processing**: Uses Go's concurrency features for parallel validation
- **RESTful API**: Clean API interface for integration with front-end applications
- **Graceful Shutdown**: Proper resource cleanup and request completion on shutdown
- **Configurable**: Easily configure server settings, database options, and performance parameters

## Tech Stack

- **Go 1.24**: Leverages the performance and concurrency features of Go
- **SQLite**: Lightweight, file-based database for persistent storage
- **gorilla/mux**: Fast and flexible HTTP router for REST endpoints
- **go-playground/validator**: Request validation
- **Context Support**: For proper request cancellation and timeouts

## Getting Started

### Prerequisites

- Go 1.24.2 or higher
- Git


### Prerequisites

- Go 1.24
- Git

### Setup

1. Clone the repository:
   ```
   git clone https://github.com/jilani-go/glofox.git
   cd glofox
   ```

2. Install dependencies:
   ```
   go mod download
   ```

3. Build and run the application:
   ```
   make build
   make run
   ```

## API Documentation

The API is documented using OpenAPI Specification (OAS) which provides a standardized way to describe RESTful APIs. 
