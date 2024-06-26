#!/bin/bash

# Set the default port if not set
PORT=${PORT:-8418}

# Start the Go API server
./main &

# Start the Next.js server
cd /app
npm run start -- --port $PORT
