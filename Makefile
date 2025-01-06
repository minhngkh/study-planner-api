binary_name = main
lambda_name = lambda


## all: build & test the application
all: build test

build: 
	@echo "Building..."
	
	
	@go build -o main cmd/server/main.go

## run: run the application
run:
	@go run cmd/server/main.go

## gen
gen:
	@go generate ./...

## gen/models: generate model schemas using gorm-gen
gen/models:
	@go generate ./internal/model

## gen/api: generate api files using oapi-codegen
gen/api:
	@go generate ./internal/api

## test: run tests
test:
	@echo "Testing..."
	@go test ./... -v

## clean
clean:
	@echo "Cleaning..."
	@rm -rf  tmp $(binary_name) $(lambda_name).zip

## watch: run the application with air
watch:
	@go run github.com/air-verse/air

## build/lambda: build binary for aws lambda
build/lambda:
	@echo "Building Lambda..."
	@GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o bootstrap -tags lambda.norpc cmd/lambda/lambda.go
	@zip -r lambda.zip bootstrap .env .env.local templates
	@rm bootstrap


## help: print this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'


.PHONY: all build run test clean watch build-lambda gen-models gen-api
