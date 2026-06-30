# Variables
BINARY_NAME=my-ls
CMD_PATH=./cmd/my-ls

# Intercept command-line arguments after "make run"
ifeq ($(firstword $(MAKECMDGOALS)),run)
  RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(RUN_ARGS):;@:)
endif

.PHONY: all build clean test run help

all: clean test build

build:
	@echo "🔨 Building the binary..."
	@go build -o $(BINARY_NAME) $(CMD_PATH)
	@echo "✅ Build complete! Binary located at: ./$(BINARY_NAME)"

run:
	@go run $(CMD_PATH) $(RUN_ARGS)

test:
	@echo "🧪 Running unit tests..."
	@go test -v ./...

clean:
	@echo "🧹 Cleaning up build artifacts..."
	@rm -rf ./$(BINARY_NAME)
	@go clean
	@echo "✨ Clean complete."

help:
	@echo "Available targets:"
	@echo "  make build      - Compile the Go binary into the project root"
	@echo "  make run        - Run the application directly using go run"
	@echo "  make run --     - Usage with flags, i.e: make run -- -l -a"
	@echo "  make test       - Execute all unit tests"
	@echo "  make clean      - Remove build directories and binary files"