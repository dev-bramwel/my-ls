# Variables
BINARY_NAME=my-ls
BUILD_DIR=build
CMD_PATH=./cmd/my-ls

.PHONY: all build clean test run help

all: clean test build

build:
	@echo "🔨 Building the binary..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_PATH)
	@echo "✅ Build complete! Binary located at: $(BUILD_DIR)/$(BINARY_NAME)"

run:
	@go run $(CMD_PATH)

test:
	@echo "🧪 Running unit tests..."
	@go test -v ./...

clean:
	@echo "🧹 Cleaning up build artifacts..."
	@rm -rf $(BUILD_DIR)
	@go clean
	@echo "✨ Clean complete."

help:
	@echo "Available targets:"
	@echo "  make build  - Compile the Go binary into $(BUILD_DIR)/"
	@echo "  make run    - Run the application directly using go run"
	@echo "  make test   - Execute all unit tests"
	@echo "  make clean  - Remove build directories and binary files"