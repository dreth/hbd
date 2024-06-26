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
FROM node:20-alpine AS node-builder

# Move to working directory (/app).
WORKDIR /app

# Copy the Next.js app files
COPY frontend/package*.json ./
RUN npm install

# Copy the rest of the frontend code
COPY frontend/ ./

# Build the Next.js app
RUN npm run build

# Stage 3: Create the final container
FROM node:20-alpine AS final

# Install additional tools
RUN apk add --no-cache ca-certificates bash

# Copy the Go API binary
COPY --from=go-builder /build/main /main

# Copy the Next.js build output
COPY --from=node-builder /app/.next /app/.next
COPY --from=node-builder /app/public /app/public
COPY --from=node-builder /app/package.json /app/package.json
COPY --from=node-builder /app/package-lock.json /app/package-lock.json
COPY --from=node-builder /app/node_modules /app/node_modules

# Command to run both the API server and the Next.js server
# Using a bash script to start both services
COPY start.sh /start.sh
RUN chmod +x /start.sh

# Expose the necessary ports
EXPOSE 8418

ENTRYPOINT ["/start.sh"]
