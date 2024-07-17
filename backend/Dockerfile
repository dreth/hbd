# Stage 1: Build the Go application
FROM golang:1.22.4-alpine AS builder

# Install necessary dependencies for building the application and SQLite
RUN apk add --no-cache gcc libc-dev sqlite-dev

# Move to working directory (/build)
WORKDIR /build

# Copy and download dependencies using go mod
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY . .

# Set necessary environment variables needed for the image and build the API server
ENV CGO_ENABLED=1 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o main .

# Stage 2: Create the final lightweight image
FROM alpine:latest

# Install necessary runtime dependencies for SQLite
RUN apk add --no-cache sqlite

# Copy the binary and config files from the builder stage
COPY --from=builder /build/main /main
COPY --from=builder /build/migrations /migrations

# Set the entry point to the binary
ENTRYPOINT ["/main"]
