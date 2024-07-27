#!/bin/bash

# Set the default port if not set
PORT=${PORT:-8418}
cd /app/backend

# if 'hbd.db' does not exist, create it as an empty file
if [ ! -f /app/data/hbd.db ]; then
    touch /app/data/hbd.db
fi

# Start the Go API server
cd /app/frontend && npm run build

# Start the Next.js server
# Build the frontend and run it
/app/backend/main &
npm run start -- --port $PORT --hostname 0.0.0.0
