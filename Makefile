# Variables
BINARY_NAME=sgit
VERSION?=v1.0.0
BUILD_DIR=build
MAIN_PACKAGE=./cmd/main.go

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-ldflags "-X github.com/hunkim/sgit/cmd.version=$(VERSION)"

.PHONY: all build clean test deps tidy run install

# Default target
all: clean deps test build

# Build the binary
build:
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .

# Clean build files
clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

# Run tests
test:
	$(GOTEST) -v ./...

# Download dependencies
deps:
	$(GOMOD) download

# Tidy up dependencies
tidy:
	$(GOMOD) tidy

# Run the application
run:
	$(GOCMD) run . 

# Install locally
install:
	$(GOCMD) install $(LDFLAGS) .

# Build for multiple platforms
build-all: clean deps
	mkdir -p $(BUILD_DIR)
	
	# Linux AMD64
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
	
	# Linux ARM64
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 .
	
	# macOS AMD64 (Intel)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	
	# macOS ARM64 (Apple Silicon)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .
	
	# Windows AMD64
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .

# Create release packages
release: build-all
	cd $(BUILD_DIR) && \
	tar -czf $(BINARY_NAME)-$(VERSION)-linux-amd64.tar.gz $(BINARY_NAME)-linux-amd64 && \
	tar -czf $(BINARY_NAME)-$(VERSION)-linux-arm64.tar.gz $(BINARY_NAME)-linux-arm64 && \
	tar -czf $(BINARY_NAME)-$(VERSION)-darwin-amd64.tar.gz $(BINARY_NAME)-darwin-amd64 && \
	tar -czf $(BINARY_NAME)-$(VERSION)-darwin-arm64.tar.gz $(BINARY_NAME)-darwin-arm64 && \
	zip $(BINARY_NAME)-$(VERSION)-windows-amd64.zip $(BINARY_NAME)-windows-amd64.exe

# Development helpers
dev-setup:
	$(GOGET) -u github.com/spf13/cobra/cobra
	$(GOMOD) tidy

# Quick test build and run
quick: build
	./$(BUILD_DIR)/$(BINARY_NAME) --help 