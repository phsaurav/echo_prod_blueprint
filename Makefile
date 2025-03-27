include .envrc

AIR := $(HOME)/go/bin/air

.PHONY: run debug test migration migrate-up migrate-down migrate-down-to migration-tree seed gen-docs

all: build test

build:
	@echo "Building..."
	
	
	@go build -o main cmd/api/main.go

# Run the application
run:
	$(AIR)

debug:
	DEBUG=true $(AIR)

.PHONY: migrate-create
migration:
	@goose -s create $(filter-out $@,$(MAKECMDGOALS)) sql

.PHONY: migrate-up
migrate-up:
	@goose up

.PHONY: migrate-down
migrate-down:
	@goose down

.PHONY: migrate-down-to
migrate-down-to:
	@goose down-to $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migration-tree
migration-tree:
	@goose status$(filter-out $@,$(MAKECMDGOALS))


# Create DB container
docker-run:
	@if docker compose up --build 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up --build; \
	fi

# Shutdown DB container
docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v
# Integrations Tests for the application
itest:
	@echo "Running integration tests..."
	@go test ./internal/database -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi

.PHONY: seed
seed:
	@go run cmd/migrate/seed/main.go

.PHONY: gen-docs
gen-docs:
	@swag init -g ./api/main.go -d cmd,internal && swag fmt
