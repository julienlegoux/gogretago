# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/api

# Run stage
FROM alpine:3.21

WORKDIR /app

# Copy binary from builder
COPY --from=builder /api .

# Expose port
EXPOSE 3000

# Run the application
CMD ["./api"]
