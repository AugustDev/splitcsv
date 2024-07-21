# Binary name
BINARY_NAME=splitcsv

# Source files
SRC=./cmd

# Output directory
OUTPUT_DIR=./bin

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Build for the current OS
.PHONY: build
build:
	$(GOBUILD) -o $(OUTPUT_DIR)/$(BINARY_NAME) -v $(SRC)

# Clean build artifacts
.PHONY: clean
clean:
	$(GOCLEAN)
	rm -rf $(OUTPUT_DIR)

# Run tests
.PHONY: test
test:
	$(GOTEST) -v ./...

# Build for all target operating systems
.PHONY: build-all
build-all: clean
	mkdir -p $(OUTPUT_DIR)
	# Build for Windows
	GOOS=windows GOARCH=386 $(GOBUILD) -o $(OUTPUT_DIR)/$(BINARY_NAME)_windows_386.exe $(SRC)
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(OUTPUT_DIR)/$(BINARY_NAME)_windows_amd64.exe $(SRC)
	# Build for macOS
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(OUTPUT_DIR)/$(BINARY_NAME)_darwin_amd64 $(SRC)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(OUTPUT_DIR)/$(BINARY_NAME)_darwin_arm64 $(SRC)
	# Build for Linux
	GOOS=linux GOARCH=386 $(GOBUILD) -o $(OUTPUT_DIR)/$(BINARY_NAME)_linux_386 $(SRC)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(OUTPUT_DIR)/$(BINARY_NAME)_linux_amd64 $(SRC)
	GOOS=linux GOARCH=arm $(GOBUILD) -o $(OUTPUT_DIR)/$(BINARY_NAME)_linux_arm $(SRC)
	GOOS=linux GOARCH=arm64 $(GOBUILD) -o $(OUTPUT_DIR)/$(BINARY_NAME)_linux_arm64 $(SRC)

# Run the application (assumes it's built for the current OS)
.PHONY: run
run: build
	$(OUTPUT_DIR)/$(BINARY_NAME) $(ARGS)

# Run without building
.PHONY: dev
dev:
	$(GOCMD) run $(SRC)/main.go $(ARGS)