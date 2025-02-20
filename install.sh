#!/bin/bash

echo "Installing gomon ..."

#check if gois installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go first."
    exit 1
fi

#insatll gomon
go install github.com/srivatsa-bot/gomon
echo "Aded gomon to your local go bin"

#check if gomon is installed

if command -v gomon &> /dev/null; then
    echo "✓ gomon installed successfully!"
    echo "You can now use 'gomon' command to watch your Go files."
else
    echo "Error: Installation failed. Please make sure your GOPATH/bin is in your PATH"
    echo "Add this to your ~/.bashrc or ~/.zshrc:"
    echo "export PATH=\$PATH:\$(go env GOPATH)/bin"
fi
