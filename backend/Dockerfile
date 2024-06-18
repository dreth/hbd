# Use the same Go version for consistency
FROM golang:1.22.4-alpine AS builder

# Install git - required for fetching go dependencies.
RUN apk add --no-cache git

WORKDIR /build

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Set necessary environmet for building the Go application
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Build the Go app
RUN go build -ldflags="-s -w" -o main .

# Build migrate tool
RUN go get -tags 'postgres' -u github.com/golang-migrate/migrate/v4/cmd/migrate/

# Start a new stage from scratch
FROM scratch

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /build/main .
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate
COPY ./migrations /migrations

# Command to run
ENTRYPOINT ["/usr/local/bin/migrate", "-database", "postgres://postgres:postgres@localhost:6684/postgres?sslmode=disable", "-path", "/migrations", "up", "/main"]
