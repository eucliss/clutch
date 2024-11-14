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
	"clutch/store"
)

// func CreateIndex(index string) {
// 	newIndex := store.Index{
// 		Name: index,
// 		Mapping: `
// 		{
// 		  "settings": {
// 			"number_of_shards": 1
// 		  },
// 		  "mappings": {
// 			"properties": {
// 			  "Name": {
// 				"type": "text"
// 			  },
// 			  "Description": {
// 				"type": "text"
// 			  },
// 			  "Hostname": {
// 				"type": "text"
// 			  },
// 			  "Time": {
// 				"type": "text"
// 			  }
// 			}
// 		  }
// 		}`,
// 	}
// 	c := common.GetConfig()
// 	c.Store.CreateIndices(newIndex)
// }

// Add a new factory function
func InitializeStore(cfg *common.DatabaseConfig) (common.Store, error) {
	switch cfg.Type {
	case "elastic":
		fmt.Println("Initializing Elasticsearch...")
		store := &store.ElasticStore{
			Location: cfg.CertLocation,
			Address:  fmt.Sprintf("https://%s:%s", cfg.Host, cfg.Port),
		}
		store.SetUsername(cfg.User)
		store.SetPassword(cfg.Password)
		store.Initialize()
		return store, nil
	case "qdrant":
		fmt.Println("Initializing Qdrant...")
		store := &store.QdrantStore{
			Host: cfg.Host,
			Port: cfg.Port,
		}
		store.Initialize()
		return store, nil
	default:
		return nil, fmt.Errorf("unsupported store type: %s", cfg.Type)
	}
}

func DeleteIndex(c common.Store, index string) {
	c.DeleteIndex(index)
}

func main() {
	// Load the base config
	base_config, err := config.LoadCommonConfig()
	// Set the common config value for use across the program
	common.SetConfig(base_config)
	cfg := &common.GlobalConfig
	if err != nil {
		fmt.Println("Error loading base config:", err)
		return
	}
	fmt.Println("Common.config loaded")

	store, err := InitializeStore(&cfg.Database)
	if err != nil {
		fmt.Println("Error initializing store:", err)
		return
	}
	// Update common config with the store config
	cfg.SetStoreConfig(store)

	// return
	// Start the services
	go services.Start(&common.Pipeline)
	// Initialize the receiver to get events from websocket
	r := receiver.NewReceiver()
	// Start the receiver
	r.Receive()

	// wait for time to sleep
	fmt.Println("Sleeping for 10 seconds...")
	time.Sleep(10 * time.Second)
	fmt.Println("Continuing execution...")

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
	return

	// Start the websocket server to listen for events
	fmt.Println("Starting WebSocket server on :8080")
	if err := r.StartServer(":8080"); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
