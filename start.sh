#!/bin/bash

# Set the default port if not set
PORT=${PORT:-8418}
cd /app/backend

# if 'hbd.db' does not exist, create it as an empty file
if [ ! -f hbd.db ]; then
    touch /app/data/hbd.db
fi

# Start the Go API server
/app/backend/main &

# Start the Next.js server
# Build the frontend and run it
cd /app/frontend && npm run build
npm run start -- --port $PORT --hostname 0.0.0.0
