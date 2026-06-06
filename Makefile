.PHONY: all build build-ui build-server clean test help

# Default target
all: build

# Help target
help:
	@echo "Available targets:"
	@echo "  build        - Build the entire application (frontend + backend + cli)"
	@echo "  build-ui     - Build the React frontend"
	@echo "  build-server - Build the Go backend binary"
	@echo "  build-cli    - Build the Ingestion CLI binary"
	@echo "  test         - Run backend tests"
	@echo "  clean        - Remove build artifacts"

# Build everything
build: build-ui build-server build-cli

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

# Build the ingestion CLI
build-cli:
	@echo "Building CLI..."
	mkdir -p bin
	cd telemetry-mend-cli && go build -o ../bin/tm-ingest main.go
	@echo "Binary created at bin/tm-ingest"

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
