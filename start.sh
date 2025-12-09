#!/bin/bash
# CloudEvents Explorer - Quick Launcher

echo "🚀 Starting CloudEvents Explorer..."
echo ""
echo "📡 Opening http://localhost:8888 in your browser..."
echo "🛑 Press Ctrl+C to stop the server"
echo ""

# Open browser after a short delay
(sleep 2 && open http://localhost:8888) &

# Start the server
go run main.go