# Base image
FROM golang:1.19

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
COPY docs/order docs/order
COPY pkg/ pkg/
COPY docs/order docs/order
COPY internal/configs internal/configs
COPY internal/models internal/models
COPY internal/repository internal/repository

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN go build -o app ./internal/apps/order-api/

# Expose port 80 for the API
EXPOSE 80

# Command to run the application
ENTRYPOINT ["./app"]