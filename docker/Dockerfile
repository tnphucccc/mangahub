FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build all binaries
RUN go build -o api-server ./cmd/api-server/main.go
RUN go build -o tcp-server ./cmd/tcp-server/main.go
RUN go build -o udp-server ./cmd/udp-server/main.go
RUN go build -o grpc-server ./cmd/grpc-server/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates

# Copy binaries from builder
COPY --from=builder /app/api-server .
COPY --from=builder /app/tcp-server .
COPY --from=builder /app/udp-server .
COPY --from=builder /app/grpc-server .

# Create necessary directories
RUN mkdir -p /app/data /app/configs

# Default command
CMD ["./api-server"]
