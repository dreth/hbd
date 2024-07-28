# Stage 1: Build the Go API server
FROM --platform=linux/amd64 golang:1.22.4-alpine AS go-builder

# add gcc
RUN apk add --no-cache gcc libc-dev sqlite-dev musl-dev

# Move to working directory (/build).
WORKDIR /build

# Copy and download dependency using go mod.
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy the backend code into the container.
COPY backend/ ./

# Set necessary environment variables needed for our image and build the API server.
ENV CGO_ENABLED=1 GOOS=linux GOARCH=amd64
RUN go build -tags musl --ldflags "-extldflags -static -s -w" -o main .

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
RUN apk add --no-cache tzdata sqlite sqlite-dev ca-certificates bash gcc

# Copy the Go API binary
COPY --from=go-builder /build/main ./backend/main

# start.sh runs both the API server and the Next.js server
# Using a bash script to start both services
COPY prod.start.sh prod.start.sh
RUN chmod +x prod.start.sh

# Expose the necessary ports
EXPOSE 8418 8417

# Create the /app/data directory
RUN mkdir /app/data

ENTRYPOINT ["/app/prod.start.sh"]
