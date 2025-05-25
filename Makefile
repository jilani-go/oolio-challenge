.PHONY: build run test lint clean

# Build the application
build:
	go build -o bin/app cmd/*.go

# Run the application
run:
	go run cmd/*.go

# Run all tests
test:
	go test ./...

# Lint the code
lint:
	golangci-lint run

# Clean build artifacts
clean:
	rm -rf bin/