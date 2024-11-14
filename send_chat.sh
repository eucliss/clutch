#!/bin/bash

# WebSocket server URL
WS_URL="ws://localhost:8080/chat"

# Path to the JSON file
TXT_FILE="data/request.txt"


# Read the TXT content and remove newlines
CONTENT=$(cat "$TXT_FILE" | tr -d '\n' | tr -d '\r')

# Construct the payload
PAYLOAD="{\"type\":\"chat\",\"payload\":$CONTENT}"
# Send the payload to the WebSocket server using websocat
echo "$PAYLOAD" | websocat "$WS_URL"

echo "Payload: $PAYLOAD"

echo "Event sent successfully"
