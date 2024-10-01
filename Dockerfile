# Build stage
FROM golang:1.23-alpine AS build

WORKDIR /app

# Copy go.mod and go.sum files for dependency caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code
COPY . .

# Build the Go application
RUN go build -o proxy-app ./cmd/proxy/main.go

# Final stage
FROM alpine:3.17

WORKDIR /app

# Copy the binary and necessary files from the build stage
COPY --from=build /app/proxy-app .
COPY configs/config.yaml ./configs/
COPY cookieforce.txt ./

# Set the entry point for the container
CMD ["./proxy-app"]
