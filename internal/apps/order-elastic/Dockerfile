# base image from 1.19-alpine for minimum size
FROM golang:1.19-alpine AS builder

# set working directory
WORKDIR /app

# install gcc
RUN apk add --no-cache gcc musl-dev

# copy all files
COPY . .

## cache deps
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

## build the binary
WORKDIR /app
RUN CGO_ENABLED=0
RUN GOOS=linux
RUN GO111MODULE=on
RUN GOARCH=amd64
RUN go build -ldflags="-w -s" -tags musl -o main .

# stage-2: image builder
FROM alpine
WORKDIR /build
ENV project="orderElastic"
COPY --from=builder /app/main .

# run
RUN chmod +x /build/main
ENTRYPOINT [ "/build/main" ]