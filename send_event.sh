#!/bin/bash

# WebSocket server URL
WS_URL="ws://localhost:8080/ws"

# Path to the JSON file
JSON_FILE="data/machinery.json"

# Check if the JSON file exists
if [ ! -f "$JSON_FILE" ]; then
    echo "Error: JSON file not found at $JSON_FILE"
    exit 1
fi

# Read the JSON content and remove newlines
JSON_CONTENT=$(cat "$JSON_FILE" | tr -d '\n' | tr -d '\r')

# Construct the payload
PAYLOAD="{\"type\":\"new_collection_testing\",\"payload\":$JSON_CONTENT}"
# Send the payload to the WebSocket server using websocat
echo "$PAYLOAD" | websocat "$WS_URL"

echo "Payload: $PAYLOAD"

echo "Event sent successfully"
