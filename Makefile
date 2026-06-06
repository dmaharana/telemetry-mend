.PHONY: all build build-ui build-server clean test help

# Default target
all: build

# Help target
help:
	@echo "Available targets:"
	@echo "  build        - Build the entire application (frontend + backend)"
	@echo "  build-ui     - Build the React frontend"
	@echo "  build-server - Build the Go backend binary"
	@echo "  test         - Run backend tests"
	@echo "  clean        - Remove build artifacts"

# Build everything
build: build-ui build-server

# Build the frontend
build-ui:
	@echo "Building frontend..."
	cd telemetry-mend-ui && pnpm install && pnpm run build
	@echo "Embedding frontend assets..."
	mkdir -p telemetry-mend-server/internal/web/dist
	cp -r telemetry-mend-ui/dist/* telemetry-mend-server/internal/web/dist/

# Build the backend server
build-server:
	@echo "Building backend..."
	mkdir -p bin
	cd telemetry-mend-server && go build -o ../bin/telemetry-mend ./cmd/server/main.go
	@echo "Binary created at bin/telemetry-mend"

# Run tests
test:
	@echo "Running tests..."
	cd telemetry-mend-server && go test ./...

# Clean build artifacts
clean:
	@echo "Cleaning artifacts..."
	rm -rf bin
	rm -rf telemetry-mend-ui/dist
	rm -rf telemetry-mend-server/internal/web/dist
	rm -rf telemetry-mend-server/telemetry_mend.db
