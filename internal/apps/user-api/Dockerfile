# base image from 1.19-alpine for minimum size
FROM golang:1.19-alpine AS builder

# set working directory
WORKDIR /app

# install gcc
RUN apk add --no-cache gcc musl-dev

# copy all files
COPY . .

# specificly install for swagger this version
RUN go install github.com/swaggo/swag/cmd/swag@v1.8.3

## cache deps
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

## build the binary
WORKDIR /app
RUN swag init --g ./cmd/order-api.go --o ./docs/order -instanceName orderAPI --exclude ./internal/apps/user-api
RUN CGO_ENABLED=0
RUN GOOS=linux
RUN GO111MODULE=on
RUN GOARCH=amd64
RUN go build -ldflags="-w -s" -tags musl -o main .

# stage-2: image builder
FROM alpine
WORKDIR /build
ENV project="orderAPI"
COPY --from=builder /app/main .
COPY --from=builder /app/docs .

# run
RUN chmod +x /build/main
ENTRYPOINT [ "/build/main" ]