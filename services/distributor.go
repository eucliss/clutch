package services

import (
	"clutch/common"
	"clutch/model"
	"clutch/services/mask"
	"clutch/services/storage"
	"clutch/services/synth"
	"encoding/json"
	"fmt"
)

func InitializeModel() {
	cfg := common.GetConfigAddress()
	model, err := model.NewModel(cfg.Model.URL, cfg.Model.ModelName)
	if err != nil {
		fmt.Println("Error creating model:", err)
	}

	cfg.SetModelConfig(model)
	fmt.Println("Model initialized:", model)

	// Testing the model quick
	var c = map[string]interface{}{
		"machine_id": "4",
	}
	flattenedBody := make(map[string]interface{})
	common.FlattenMap(flattenedBody, "", c)

	// Convert document to string for embedding
	jsonStr, err := json.Marshal(flattenedBody)
	q := fmt.Sprintf("What field is this machine in? %s", string(jsonStr))
	fmt.Println("Question:", q)
	res, err := model.Ask(string(jsonStr))
	if err != nil {
		fmt.Println("Error asking model:", err)
	}
	fmt.Println("Model response:", res)
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
		case "synth":
			go synth.Synth(&common.SynthChan)
		}
	}
}

func Start(pipeline *chan common.Event) {
	fmt.Println("Distributor started, priming services.")
	prime()
	for event := range *pipeline {
		fmt.Println("Distributing event:", event)
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
