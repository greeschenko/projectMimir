# project root Makefile

# List of services in monorepo
SERVICES := auth-service

# -----------------------------
# Build all services binaries
# -----------------------------
.PHONY: build
build: build-all docker-build up
	@echo "Build finished for all services."

# Build binaries for all services
.PHONY: build-all
build-all:
	@echo "Building binaries for all services..."
	@for service in $(SERVICES); do \
		echo "Building $$service..."; \
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $$service/app ./$$service/cmd/main.go; \
	done

# -----------------------------
# Build Docker images for all services
# -----------------------------
.PHONY: docker-build
docker-build:
	@echo "Building Docker images for all services..."
	@for service in $(SERVICES); do \
		echo "Building Docker image for $$service..."; \
		docker build -f $$service/Dockerfile -t projectmimir-$$service .; \
	done

# -----------------------------
# Run all services via docker-compose
# -----------------------------
.PHONY: up
up:
	@echo "Starting all services via docker-compose..."
	docker-compose up

# -----------------------------
# Clean all build artifacts
# -----------------------------
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@for service in $(SERVICES); do \
		rm -rf $$service/app $$service/build; \
	done
