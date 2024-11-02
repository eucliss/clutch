package storage

import (
	"clutch/common"
	"fmt"
)

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
