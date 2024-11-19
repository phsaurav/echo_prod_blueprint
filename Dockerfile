# The build stage
FROM golang:1.23.3 as builder
WORKDIR /application
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api cmd/api/*.go

# The run stage
FROM scratch
WORKDIR /application
# Copy CA certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /application/api .
EXPOSE 8080
CMD ["./api"]