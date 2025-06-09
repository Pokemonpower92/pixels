# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go mod files (go.sum is optional)
COPY go.mod ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the application
ARG SERVICE_NAME
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/service ./cmd/$SERVICE_NAME

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy with fixed name
COPY --from=builder /app/service /app/service

# Expose the port your app runs on
EXPOSE 8080

# Run the binary
CMD ["/app/service"]