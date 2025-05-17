FROM golang:1.24-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main cmd/api/main.go

FROM alpine:3.20.1 AS prod

RUN apk --no-cache add ca-certificates tzdata && \
    mkdir -p /app/migrations

WORKDIR /app
COPY --from=build /app/main /app/main
COPY --from=build /app/.env /app/.env
COPY --from=build /app/migrations /app/migrations

RUN chmod +x /app/main

# Create a non-root user and set proper permissions
RUN adduser -D appuser && \
    chown -R appuser:appuser /app

USER appuser

# The port will be set by the environment variable at runtime
EXPOSE 8080

CMD ["/bin/sh", "-c", "echo 'Starting application...' && /app/main 2>&1"]