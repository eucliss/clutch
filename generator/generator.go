package generator

import (
	"fmt"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

func Connect() {
	// Connect to the websocket server
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		fmt.Println("Error connecting to WebSocket server:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to WebSocket server")

	// Keep the connection alive
	go func() {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("Generator reading message error:", err)
				return
			}
		}
	}()
}

func SendEvent(event string) {
	// Connect to the websocket server
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		fmt.Println("Error connecting to WebSocket server:", err)
		return
	}
	defer conn.Close()

	// Create the event payload
	payload := map[string]interface{}{
		"type":    "event",
		"payload": event,
	}

	// Send the event
	err = conn.WriteJSON(payload)
	if err != nil {
		fmt.Println("Error sending event:", err)
		return
	}

	fmt.Println("Event sent successfully:", event)
}

func EmitEvents(interval int) {

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()
	fmt.Println("Emitting events every", interval, "seconds")

	// read in the data/machinery.json file
	data, err := os.ReadFile("data/machinery.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// send the event
	SendEvent(string(data))
}
