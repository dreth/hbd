# Stage 1: Build the Go API server
FROM golang:1.22.4-alpine AS go-builder

# Move to working directory (/build).
WORKDIR /build

# Copy and download dependency using go mod.
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy the backend code into the container.
COPY backend/ ./

# Set necessary environment variables needed for our image and build the API server.
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o main .

# Stage 2: Build the Next.js app
FROM node:20-alpine

# Move to working directory (/app).
WORKDIR /app

# Copy the Next.js app files
RUN mkdir frontend
COPY frontend/package*.json ./frontend/
RUN cd frontend && npm install

# Copy the rest of the frontend code
COPY frontend/ ./frontend

# Copy the backend code into the container.
COPY backend/ ./backend

# Build the Next.js app
RUN cd frontend && npm run build

# Install additional tools
RUN apk add --no-cache tzdata sqlite ca-certificates bash

# Copy the Go API binary
COPY --from=go-builder /build/main ./backend/main

# Command to run both the API server and the Next.js server
# Using a bash script to start both services
COPY start.sh ./start.sh
RUN chmod +x ./start.sh

# Expose the necessary ports
EXPOSE 8418 8417

ENTRYPOINT ["/app/start.sh"]
