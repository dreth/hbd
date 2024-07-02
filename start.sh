#!/bin/bash

# Set the default port if not set
PORT=${PORT:-8418}
cd /app/backend

# Start the Go API server
./main &

# Start the Next.js server
cd ../frontend
npm run start -- --port $PORT
