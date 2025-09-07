#!/bin/bash

# Run gosec security scanner excluding auto-generated protobuf files
echo "Running gosec security scan (excluding auto-generated files)..."

# First verify the code compiles
echo "Verifying code compiles..."
if ! go build -v ./...; then
    echo "❌ Code compilation failed. Please fix compilation errors first."
    exit 1
fi
echo "✅ Code compilation successful"

# Run gosec
/home/gurkanindibay/go/bin/gosec \
  -exclude-dir=proto/auth \
  -exclude-dir=proto/event \
  -fmt=text \
  -out=gosec-report.txt \
  ./...

echo "Gosec scan completed. Report saved to gosec-report.txt"