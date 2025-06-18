FROM golang:alpine as builder

# Install git for dependencies
RUN apk add --no-cache git

WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the application
RUN go build -o main .
RUN ls -la

# Stage 2: Create a minimal image
FROM alpine:latest

# Add necessary runtime dependencies
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main /app/main
COPY --from=builder /app/.env /app/.env

# Verify the file exists
RUN ls -la /app

EXPOSE 3000

CMD ["/app/main"]
