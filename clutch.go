package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"clutch/common"
	"clutch/config"
	"clutch/receiver"
	"clutch/services"
	"clutch/services/storage"
)

func Start(websocket bool) receiver.Receiver {
	// Load the base config
	success, err := config.InitializeConfig()
	if !success {
		fmt.Println("Error initializing config:", err)
		return receiver.Receiver{}
	}

	// Point to the global config for future steps
	cfg := &common.GlobalConfig
	fmt.Println("Common.GlobalConfig loaded")

	store, err := storage.InitializeStore(&cfg.Database)
	if err != nil {
		fmt.Println("Error initializing store:", err)
		return receiver.Receiver{}
	}
	// Update common config with the store config
	cfg.SetStoreConfig(store)

	// Start the services
	go services.Start(&common.Pipeline)

	r := receiver.NewReceiver()
	// Start the receiver
	r.Receive()

	if websocket {
		fmt.Println("Starting WebSocket server on :8080")
		if err := r.StartServer(":8080"); err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}
	return *r
}

func main() {
	// Start the receiver (websocket / chat)
	START_WEBSOCKET := false
	reciever := Start(START_WEBSOCKET)

	fmt.Println("Reciever:", &reciever)
	fmt.Println("Services are up and running...")
	// sleep for 10 seconds for things to
	time.Sleep(5 * time.Second)

	// Reload config for stuff
	cfg := &common.GlobalConfig
	store := cfg.Store

	// ------------------ TESTING LOGIC ------------------

	// Perform a query to test the connection
	fmt.Println("Querying the DB")
	jsonQuery := map[string]interface{}{
		"timestamp":    "2024-01-01T00:00:00Z",
		"machine_id":   "4",
		"machine_type": "planter",
		"location":     "field_1",
		"status":       "running",
	}
	jsonQueryString, err := json.Marshal(jsonQuery)
	if err != nil {
		fmt.Println("Error marshalling jsonQuery:", err)
		return
	}
	fmt.Println("JSON Query:", string(jsonQueryString))
	docs := store.Query("clutch_testing_events", string(jsonQueryString))
	fmt.Println("Query results:", docs)

	// // Start the WebSocket here if you need to test events coming in and shit manually
	// fmt.Println("Starting WebSocket server on :8080")
	// if err := reciever.StartServer(":8080"); err != nil {
	// 	log.Fatal("ListenAndServe: ", err)
	// }

	return
}

// func DeleteIndex(c common.Store, index string) {
// 	c.DeleteIndex(index)
// }
