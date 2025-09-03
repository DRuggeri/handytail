# termtail Makefile

# Variables
BINARY_NAME=termtail
BINARY_PATH=./$(BINARY_NAME)
GO_FILES=$(shell find . -name "*.go" -type f)
VERSION?=dev
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

# Default target
.PHONY: all
all: clean build test

# Build the binary
.PHONY: build
build: $(BINARY_NAME)

$(BINARY_NAME): $(GO_FILES)
	go build $(LDFLAGS) -o $(BINARY_NAME) .

# Run tests
.PHONY: test
test:
	go test -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run benchmarks
.PHONY: bench
bench:
	go test -bench=. -benchmem ./...

# Clean build artifacts
.PHONY: clean
clean:
	go clean
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html

# Install binary to GOPATH/bin
.PHONY: install
install: build
	go install $(LDFLAGS) .

# Run linting
.PHONY: lint
lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed, running go vet instead"; \
		go vet ./...; \
	fi

# Format code
.PHONY: fmt
fmt:
	go fmt ./...

# Check formatting
.PHONY: fmt-check
fmt-check:
	@if [ -n "$$(go fmt ./...)" ]; then \
		echo "Code is not formatted. Run 'make fmt' to fix."; \
		exit 1; \
	fi

# Run all checks (format, lint, test)
.PHONY: check
check: fmt-check lint test

# Build for multiple platforms
.PHONY: build-all
build-all: clean
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BINARY_NAME)-linux-arm64 .
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BINARY_NAME)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-windows-amd64.exe .

# Create release archives
.PHONY: release
release: build-all
	@mkdir -p dist
	@for binary in $(BINARY_NAME)-*; do \
		if [ -f "$$binary" ]; then \
			echo "Creating archive for $$binary"; \
			tar -czf "dist/$$binary.tar.gz" "$$binary" README.md; \
		fi \
	done
	@echo "Release archives created in dist/"

# Development setup
.PHONY: dev-setup
dev-setup:
	go mod tidy
	@echo "Installing development tools..."
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi

# Quick test run with example
.PHONY: example
example: build
	@echo "Running example test..."
	@echo -e "Processing...\nBUILD SUCCESS\nDone" | ./$(BINARY_NAME) -success "SUCCESS" -failure "ERROR"
	@echo "Exit code: $$?"

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all          - Clean, build, and test"
	@echo "  build        - Build the binary"
	@echo "  test         - Run tests"
	@echo "  test-coverage- Run tests with coverage report"
	@echo "  bench        - Run benchmarks"
	@echo "  clean        - Remove build artifacts"
	@echo "  install      - Install binary to GOPATH/bin"
	@echo "  lint         - Run linting"
	@echo "  fmt          - Format code"
	@echo "  fmt-check    - Check if code is formatted"
	@echo "  check        - Run format check, lint, and tests"
	@echo "  build-all    - Build for multiple platforms"
	@echo "  release      - Create release archives"
	@echo "  dev-setup    - Set up development environment"
	@echo "  example      - Run a quick example"
	@echo "  help         - Show this help message"
