# THIS IS THE LOCAL DOCKERFILE
# This Dockerfile is used to build the Docker image for local development or local personal use.
# It is not intended for production use.

# Stage 1: Build the Go API server
FROM --platform=linux/amd64 golang:1.22.4-alpine AS go-builder

# Add gcc
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

# Install Nginx
RUN apk add --no-cache nginx

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
RUN cd frontend && rm -rf .next && npm run build

# Install additional tools
RUN apk add --no-cache tzdata sqlite sqlite-dev ca-certificates bash gcc

# Copy the Go API binary
COPY --from=go-builder /build/main ./backend/main

# Copy Nginx configuration file
COPY nginx.conf /etc/nginx/nginx.conf

# Install s6 overlay
ENV S6_OVERLAY_VERSION=3.2.0.0
ADD https://github.com/just-containers/s6-overlay/releases/download/v${S6_OVERLAY_VERSION}/s6-overlay-noarch.tar.xz /tmp

# Download the correct s6 overlay based on the architecture
RUN ARCH=$(uname -m) && \
    wget -O /tmp/s6-overlay-${ARCH}.tar.xz https://github.com/just-containers/s6-overlay/releases/download/v${S6_OVERLAY_VERSION}/s6-overlay-${ARCH}.tar.xz && \
    tar -C / -Jxpf /tmp/s6-overlay-noarch.tar.xz && \
    tar -C / -Jxpf /tmp/s6-overlay-${ARCH}.tar.xz

# Expose the necessary ports
EXPOSE 80

# Create the /app/data directory
RUN mkdir /app/data

# Create service directories
RUN mkdir -p /etc/services.d/backend /etc/services.d/frontend /etc/services.d/nginx

# Copy service scripts
COPY services/backend/run /etc/services.d/backend/run
COPY services/frontend/run /etc/services.d/frontend/run
COPY services/nginx/run /etc/services.d/nginx/run

# Copy frontend finish script
COPY services/frontend/finish /etc/services.d/frontend/finish

# Make the scripts executable
RUN chmod +x /etc/services.d/backend/run /etc/services.d/frontend/run /etc/services.d/nginx/run /etc/services.d/frontend/finish

# Set s6 as the entrypoint
ENTRYPOINT ["/init"]
