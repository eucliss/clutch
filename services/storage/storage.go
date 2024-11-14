package storage

import (
	"clutch/common"
	"clutch/store"
	"fmt"
)

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
func StoreMasks(maskedStorageChan *chan common.Event) {

	fmt.Println("Starting masked storage service")
	cfg := common.GetConfig()
	store := cfg.Store

	for event := range *maskedStorageChan {

		fmt.Println("---------- Storing ----------")
		fmt.Println("New Masked or Synthesized event:", event)
		store.InsertDocument(event.Type, event.Payload)
		fmt.Println("---------- Done Storing ----------")
	}
}

func Store(storageChan *chan common.Event) {
	fmt.Println("Starting storage service")
	cfg := common.GetConfig()
	store := cfg.Store

	for event := range *storageChan {
		fmt.Println("Storing event:", event)
		store.InsertDocument(event.Type, event.Payload)
	}
}
