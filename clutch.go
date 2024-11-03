package main

import (
	"fmt"
	"log"

	"clutch/common"
	"clutch/config"
	"clutch/receiver"
	"clutch/services"
	"clutch/store"
)

func CreateIndex(index string) {
	newIndex := store.Index{
		Name: index,
		Mapping: `
		{
		  "settings": {
			"number_of_shards": 1
		  },
		  "mappings": {
			"properties": {
			  "Name": {
				"type": "text"
			  },
			  "Description": {
				"type": "text"
			  },
			  "Hostname": {
				"type": "text"
			  },
			  "Time": {
				"type": "text"
			  }
			}
		  }
		}`,
	}
	c := common.GetConfig()
	c.Store.CreateIndices(newIndex)
}

// Add a new factory function
func InitializeStore(cfg *common.DatabaseConfig) (store.Store, error) {
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

func InitializeElastic(def common.Config) (c store.ElasticStore, res *store.ElasticStore) {
	fmt.Println("Initializing Elasticsearch...")

	c = store.ElasticStore{
		Location: def.Database.CertLocation,
		Address:  fmt.Sprintf("https://%s:%s", def.Database.Host, def.Database.Port),
	}
	c.SetUsername(def.Database.User)
	c.SetPassword(def.Database.Password)
	c.Initialize()

	res = &c
	return
}

func DeleteIndex(c store.Store, index string) {
	c.DeleteIndex(index)
}

func InitializeQdrant(def common.Config) (c store.QdrantStore, res *store.QdrantStore) {
	fmt.Println("Initializing Qdrant...")

	c = store.QdrantStore{
		Host: def.Database.Host,
		Port: def.Database.Port,
	}

	c.Initialize()

	res = &c
	return
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

	// store, _ := InitializeQdrant(*cfg)
	// fmt.Println("Store initialized", store)

	// Initialize the Elasticsearch client
	// store, _ := InitializeElastic(*cfg)
	// fmt.Println("Elasticsearch Configured")

	store, err := InitializeStore(&cfg.Database)
	if err != nil {
		fmt.Println("Error initializing store:", err)
		return
	}
	// Update common config with the store config
	cfg.SetStoreConfig(store)

	// Perform a query to test the connection
	fmt.Println("Querying the DB")
	docs := store.Query("clutch_testing_events", `
		{
			"query": {
				"match_all": {}
			}
		}
	`)
	fmt.Println(store.GetResults(docs))

	// Start the services
	go services.Distribute(&common.Pipeline)
	// Initialize the receiver to get events from websocket
	r := receiver.NewReceiver()
	// Start the receiver
	r.Receive()

	// Start the websocket server to listen for events
	fmt.Println("Starting WebSocket server on :8080")
	if err := r.StartServer(":8080"); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
