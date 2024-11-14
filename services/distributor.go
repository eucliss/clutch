package services

import (
	"clutch/common"
	"clutch/services/mask"
	"clutch/services/model"
	"clutch/services/storage"
	"fmt"
)

func InitializeModel() {
	cfg := common.GetConfigAddress()
	model, err := model.InitializeModel(&cfg.ModelConfig)
	if err != nil {
		fmt.Println("Error creating model:", err)
	}
	cfg.SetModelConfig(model)
	fmt.Println("Model initialized:", model.GetModelName())
}

func prime() {
	cfg := common.GetConfig()
	// Check if mask_storage is in the services
	mask_storage := false
	for _, service := range cfg.Services {
		fmt.Println("Checking service:", service)
		if service == "mask_storage" {
			fmt.Println("Starting masked storage service")
			go storage.StoreMasks(&common.MaskedStorageChan)
			mask_storage = true
		}
	}

	for _, service := range cfg.Services {
		fmt.Println("Starting service:", service)
		switch service {
		case "model":
			go InitializeModel()
		case "storage":
			go storage.Store(&common.StorageChan)
		case "masking":
			go mask.Mask(&common.MaskChan, mask_storage)
		}
	}
}

func Start(pipeline *chan common.Event) {
	fmt.Println("Distributor started, priming services.")
	prime()
	for event := range *pipeline {
		fmt.Println("Distributing event:", event)
		if event.Type == "chat" {
			fmt.Println("Chat event:", event)
			common.ChatChan <- event
		} else {
			for _, service := range common.GlobalConfig.Services {
				switch service {
				case "storage":
					common.StorageChan <- event
				case "masking":
					common.MaskChan <- event
				case "synth":
					common.SynthChan <- event
				}
			}
		}
	}
}
