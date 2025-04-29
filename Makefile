# Binary name
BINARY_NAME=nmeasim
BINARY_DIR=bin

# Go related variables
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/$(BINARY_DIR)
GOFILES=$(wildcard ./cmd/$(BINARY_NAME)/*.go)

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

## install: Install missing dependencies
.PHONY: install
install:
	@echo "  >  Installing dependencies..."
	go mod download
	go mod tidy

## build: Build the binary
.PHONY: build
build: clean install
	@echo "  >  Building binary..."
	go build -o $(GOBIN)/$(BINARY_NAME) ./cmd/$(BINARY_NAME)
	@echo "  >  Binary ready at $(GOBIN)/$(BINARY_NAME)"

## clean: Clean build files
.PHONY: clean
clean:
	@echo "  >  Cleaning build cache"
	go clean
	rm -f $(GOBIN)/$(BINARY_NAME)

## test: Run tests
.PHONY: test
test:
	@echo "  >  Running tests..."
	go test -v ./...

## coverage: Run tests with coverage
.PHONY: coverage
coverage:
	@echo "  >  Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "  >  Coverage report generated at coverage.html"

## fmt: Format code
.PHONY: fmt
fmt:
	@echo "  >  Formatting code..."
	go fmt ./...

## vet: Run go vet
.PHONY: vet
vet:
	@echo "  >  Running go vet..."
	go vet ./...

## lint: Run linting
.PHONY: lint
lint: fmt vet
	@echo "  >  Running linters..."
	go install golang.org/x/lint/golint@latest
	golint ./...

## run: Build and run the binary
.PHONY: run
run: build
	@echo "  >  Running binary..."
	./$(GOBIN)/$(BINARY_NAME)

## help: Show this help message
.PHONY: help
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  %-20s %s\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

# Default target
.DEFAULT_GOAL := help