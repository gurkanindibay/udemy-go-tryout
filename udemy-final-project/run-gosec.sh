#!/bin/bash

# Run gosec security scanner excluding auto-generated protobuf files
echo "Running gosec security scan (excluding auto-generated files)..."

/home/gurkanindibay/go/bin/gosec \
  -exclude-dir=proto/auth \
  -exclude-dir=proto/event \
  -fmt=text \
  -out=gosec-report.txt \
  ./...

echo "Gosec scan completed. Report saved to gosec-report.txt"