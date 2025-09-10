# Requirements:
#	- go

# Environment variables

# Tools
GOCILINT  := go tool github.com/golangci/golangci-lint/v2/cmd/golangci-lint

.DEFAULT_GOAL := help
.PHONY: help
help:
	@echo "Usage: make <target>"
	@echo "Targets:"
	@echo "  test - Run all tests"
	@echo "  test-v - Run tests with verbose output"
	@echo "  clean - Clean build artifacts"
	@echo "  lint - Run linter"

.PHONY: test
test: ## Run all tests
	go test ./...

.PHONY: test-v
test-v: ## Run tests with verbose output
	go test -v ./...

.PHONY: clean
clean: ## Clean build artifacts
	rm -rf bin/

.PHONY: lint
lint: ## Run linter
	CGO_ENABLED=1 ${GOCILINT} run ./...
