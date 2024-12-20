# Simple Makefile for a Go project

# Build the application
all: build test

build:
	@echo "Building..."
	
	
	@go build -o main cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go

# Generate models from the database
gen-models:
	@go run cmd/genmodels/main.go

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
			air; \
			echo "Watching...";\
		else \
			read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
			if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
				go install github.com/air-verse/air@latest; \
				air; \
				echo "Watching...";\
			else \
				echo "You chose not to install air. Exiting..."; \
				exit 1; \
			fi; \
		fi

# Build the lambda
build-lambda:
	@echo "Building Lambda..."
	@GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o bootstrap -tags lambda.norpc cmd/lambda/lambda.go
	@zip lambda.zip bootstrap .env
	@rm bootstrap


.PHONY: all build run test clean watch build-lambda
