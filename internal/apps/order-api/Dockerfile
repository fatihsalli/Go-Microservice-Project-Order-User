# Base image
FROM golang:1.19

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY cmd/order-api.go cmd/order-api.go
COPY go.mod go.sum ./
COPY docs/order docs/order
COPY pkg/ pkg/
COPY docs/order docs/order
COPY internal/apps/order-api internal/apps/order-api
COPY internal/configs internal/configs
COPY internal/models internal/models
COPY internal/repository internal/repository

# Download dependencies
RUN go mod download

# Build the app
RUN go build -o order-api cmd/order-api.go

# Expose port 80 for the API
EXPOSE 80

CMD ["./order-api"]