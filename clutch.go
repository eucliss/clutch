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

func InitializeElastic(def common.Config) (c store.StoreConfig, res *store.StoreConfig) {
	fmt.Println("Initializing Elasticsearch...")

	c = store.StoreConfig{
		Location: def.Database.CertLocation,
		Address:  fmt.Sprintf("https://%s:%s", def.Database.Host, def.Database.Port),
	}
	c.SetUsername(def.Database.User)
	c.SetPassword(def.Database.Password)
	c.Initialize()

	res = &c
	return
}

func DeleteIndex(c store.StoreConfig, index string) {
	c.DeleteIndex(index)
}

func main() {

	// Load the base config
	base_config, err := config.LoadCommonConfig()
	common.SetConfig(base_config)
	cfg := &common.GlobalConfig
	if err != nil {
		fmt.Println("Error loading base config:", err)
		return
	}
	fmt.Println("Common.config loaded")

	// Initialize the Elasticsearch client
	store, _ := InitializeElastic(*cfg)
	fmt.Println("Elasticsearch Configured")

	cfg.SetStoreConfig(store)

	// fmt.Println("Querying the DB")
	// docs := store.Query("clutch_testing_events", `
	// 	{
	// 		"query": {
	// 			"match_all": {}
	// 		}
	// 	}
	// `)

	// fmt.Println(store.GetResults(docs))
	go services.Distribute(&common.Pipeline)
	r := receiver.NewReceiver()
	r.Receive()

	// select {}

	fmt.Println("Starting WebSocket server on :8080")
	if err := r.StartServer(":8080"); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
