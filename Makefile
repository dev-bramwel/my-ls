# Variables
BINARY_NAME=my-ls
CMD_PATH=./cmd/my-ls

.PHONY: all build clean test run help

all: clean test build

build:
	@echo "🔨 Building the binary..."
	@go build -o $(BINARY_NAME) $(CMD_PATH)
	@echo "✅ Build complete! Binary located at: ./$(BINARY_NAME)"

run:
	@go run $(CMD_PATH) $(ARGS)

test:
	@echo "🧪 Running unit tests..."
	@go test -v ./...

clean:
	@echo "🧹 Cleaning up build artifacts..."
	@rm -rf ./$(BINARY_NAME) my_symlink target_file.txt my_dir_symlink target_dir
	@go clean
	@echo "✨ Clean complete."

help:
	@echo "Available targets:"
	@echo "  make build      - Compile the Go binary into the project root"
	@echo "  make run        - Run the application directly using go run"
	@echo "  make run ARGS="-l -a README.md" - Run with specific flags and paths"
	@echo "  make test       - Execute all unit tests"
	@echo "  make clean      - Remove build directories and binary files"