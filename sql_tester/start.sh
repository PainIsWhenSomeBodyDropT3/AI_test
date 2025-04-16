#!/bin/bash

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed or not in the PATH"
    exit 1
fi

# Install required dependencies
echo "Installing dependencies..."
go get github.com/mattn/go-sqlite3

# Run the application with web interface
echo "Starting SQL Tester with web interface..."
go run *.go -web

# This script will not reach here unless there's an error, as the web server will keep running
echo "Application exited unexpectedly." 