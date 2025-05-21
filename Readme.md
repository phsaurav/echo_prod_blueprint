# Project JonoMot

![Golang_Arch](https://github.com/user-attachments/assets/20fe6ce8-9a5a-4bb7-8d6a-18a26da4b950)

As a Software Engineer with a background in traditional OOP languages and Clean Architecture principles, I created JonoMot to bridge the gap between SOLID principles and Go's idiomatic approach to software design.

After reading "100 Go Mistakes and How to Avoid Them," I recognized the need for a reference implementation that applies these lessons adhereing SOLID principal to create maintainable, testable, and scalable Go applications. JonoMot serves as this blueprint, demonstrating how to properly structure Go APIs while following both SOLID principles and Go-specific best practices.

## SOLID Principles in JonoMot

JonoMot adheres to SOLID principles while respecting Go's philosophy:

1. **Single Responsibility Principle**: Each package has a clear, focused responsibility:
   - `internal/user`: Manages user-related functionality
   - `internal/poll`: Handles poll functionality
   - `pkg/logger`: Centralized logging
   - `pkg/response`: Standardized API responses

2. **Open/Closed Principle**: The codebase is designed for extension without modification:
   - Service interfaces like `UserService` and `PollService` can have multiple implementations
   - New features can be added via the tool in `cmd/tools/feature`

3. **Liskov Substitution Principle**: Interfaces are properly designed for substitution:
   - Repository interfaces (`user.Repository`, `poll.Repository`) can be swapped with test mocks
   - `database.Service` abstracts database operations

4. **Interface Segregation Principle**: Interfaces are client-specific and focused:
   - `PollService` defines only methods required for poll operations
   - `UserService` defines only methods needed for user functionality

5. **Dependency Inversion Principle**: Dependencies flow toward abstractions:
   - Service constructors like `NewService` accept Repository interfaces
   - `RegisterRoutes` accepts service interfaces

## Go Best Practices from "100 Go Mistakes and How to Avoid Them"

JonoMot implements numerous best practices:

1. **Code Organization**:
   - Proper package structure separating `cmd`, `internal`, and `pkg` directories
   - Clean separation between application layers (model, repository, service, routes)

2. **Error Handling**:
   - Custom error handling with `pkg/error`
   - Centralized error responses via `response.ErrorBuilder`

3. **Interfaces**:
   - Interfaces declared by consumers, not implementers
   - Focused interfaces with minimal method sets

4. **Database Handling**:
   - Connection pooling with configured MaxOpenConns/MaxIdleConns
   - Proper context usage in database operations
   - Parameterized SQL queries to prevent SQL injection

5. **Testing**:
   - Comprehensive unit tests with mocks
   - SQL mocking using sqlmock
   - Integration tests for database operations

6. **Concurrency**:
   - Graceful shutdown implementation
   - Proper context usage for cancellation

7. **Configuration**:
   - Environment-based configuration in `config`
   - Sensible defaults with overrides

8. **API Design**:
   - RESTful API with consistent response formats
   - Structured logging via `pkg/logger`
   - Swagger documentation

## Setting Up Development Environment

### Prerequisites

- Go 1.19+
- PostgreSQL or Docker
- Make

### Environment Setup

1. Clone the repository
2. Copy `.env.example` to `.env` and adjust values
3. Run database: `make docker-run`
4. Apply migrations: `make migrate-up`
5. Build and run: `make run`

### Make Commands

```bash
# Database Management
make migration create_poll    # Create a new migration for "poll" feature
make migrate-up               # Apply all pending migrations
make migrate-down             # Rollback last migration
make docker-run               # Start PostgreSQL container
make docker-down              # Stop PostgreSQL container

# Development
make build                    # Build the application
make run                      # Run the application
make watch                    # Run with live reload

# Testing
make test                     # Run unit tests
make itest                    # Run integration tests
make all                      # Build and test

# Utilities
make clean                    # Clean build artifacts
```

### API Documentation
Swagger documentation is available at `docs` when the server is running.
