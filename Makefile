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
	@echo "Creating migration: $(filter-out $@,$(MAKECMDGOALS))"
	@if [ -f .env ]; then \
			export $$(grep -v '^#' .env | xargs); \
			export GOOSE_DRIVER=postgres; \
			export GOOSE_DBSTRING="user=$$DB_USERNAME password=$$DB_PASSWORD dbname=$$DB_DATABASE host=$$DB_HOST sslmode=disable"; \
			goose -dir="$$GOOSE_MIGRATION_DIR" -s create $(filter-out $@,$(MAKECMDGOALS)) sql; \
	else \
			echo "Error: .env file not found"; \
			exit 1; \
	fi

.PHONY: migrate-up
migrate-up:
	@echo "Running migrations up..."
	@if [ -f .env ]; then \
			export $$(grep -v '^#' .env | xargs); \
			export GOOSE_DRIVER=postgres; \
			export GOOSE_DBSTRING="user=$$DB_USERNAME password=$$DB_PASSWORD dbname=$$DB_DATABASE host=$$DB_HOST sslmode=disable"; \
			goose -dir="$$GOOSE_MIGRATION_DIR" up; \
	else \
			echo "Error: .env file not found"; \
			exit 1; \
	fi

.PHONY: migrate-down
migrate-down:
	@echo "Rolling back the last migration..."
	@if [ -f .env ]; then \
			export $$(grep -v '^#' .env | xargs); \
			export GOOSE_DRIVER=postgres; \
			export GOOSE_DBSTRING="user=$$DB_USERNAME password=$$DB_PASSWORD dbname=$$DB_DATABASE host=$$DB_HOST sslmode=disable"; \
			goose -dir="$$GOOSE_MIGRATION_DIR" down; \
	else \
			echo "Error: .env file not found"; \
			exit 1; \
	fi

.PHONY: migrate-down-to
migrate-down-to:
	@echo "Rolling back to version $(filter-out $@,$(MAKECMDGOALS))..."
	@if [ -f .env ]; then \
			export $$(grep -v '^#' .env | xargs); \
			export GOOSE_DRIVER=postgres; \
			export GOOSE_DBSTRING="user=$$DB_USERNAME password=$$DB_PASSWORD dbname=$$DB_DATABASE host=$$DB_HOST sslmode=disable"; \
			goose -dir="$$GOOSE_MIGRATION_DIR" down-to $(filter-out $@,$(MAKECMDGOALS)); \
	else \
			echo "Error: .env file not found"; \
			exit 1; \
	fi

.PHONY: migration-status
migration-status:
	@echo "Migration status..."
	@if [ -f .env ]; then \
			export $$(grep -v '^#' .env | xargs); \
			export GOOSE_DRIVER=postgres; \
			export GOOSE_DBSTRING="user=$$DB_USERNAME password=$$DB_PASSWORD dbname=$$DB_DATABASE host=$$DB_HOST sslmode=disable"; \
			goose -dir="$$GOOSE_MIGRATION_DIR" status; \
	else \
			echo "Error: .env file not found"; \
			exit 1; \
	fi

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
	@swag init -g api/main.go -d ./cmd,./internal,./pkg --parseDependency --parseInternal

.PHONY: gen-feature
gen-feature:
	@echo "Generating new feature: $(filter-out $@,$(MAKECMDGOALS))"
	@mkdir -p internal/$(filter-out $@,$(MAKECMDGOALS))
	@go run cmd/tools/feature/create.go $(filter-out $@,$(MAKECMDGOALS))