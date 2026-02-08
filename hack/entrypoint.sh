#!/bin/bash
set -e

echo "-----------------------------------------------------"
echo "🚀 Volcano-Global Lab is Live!"
echo "Connected to Mac Docker via socket."
echo "-----------------------------------------------------"

# Verify tools are working
go version
kind version
docker version --format 'Client Version: {{.Client.Version}} | Server Version: {{.Server.Version}}'

echo ""
echo "Checking existing clusters..."
kind get clusters || echo "No clusters found."

echo "-----------------------------------------------------"
echo "You are now inside the container workspace."
echo "-----------------------------------------------------"

exec /bin/bash
