#!/bin/bash
# Testing Studio (CloudEvents Explorer) - Quick Launcher

echo "🚀 Starting Testing Studio..."
echo ""
echo "📡 Opening http://localhost:8888 in your browser..."
echo "🛑 Press Ctrl+C to stop the server"
echo ""

# Open browser after a short delay
(sleep 2 && open http://localhost:8888) &

# Start the server from new location
go run cmd/server/main.go