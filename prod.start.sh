#!/bin/bash

# Set the default port if not set
PORT=${PORT:-8418}
cd /app/backend

# if 'hbd.db' does not exist, create it as an empty file
if [ ! -d /app/data ]; then
    mkdir -p /app/data
fi

if [ ! -f /app/data/hbd.db ]; then
    touch /app/data/hbd.db
fi

# Start the Go API server
./main &

# Start the Next.js server
cd ../frontend
rm -rf .next
npm run start -- --port $PORT --hostname 0.0.0.0
