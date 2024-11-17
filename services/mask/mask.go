package mask

import (
	"clutch/common"
	"clutch/config"
	"clutch/services/operations"
	"fmt"
	"path/filepath"
	"strings"
)

var MaskedEvents []MaskedEvent

var MaskMap = make(map[string]common.MaskConfig)

type MaskedEvent struct {
	RawEvent    common.Event
	MaskedEvent common.Event
	Type        string
}

func (m *MaskedEvent) confirmStringField(key string) bool {
	if _, ok := m.MaskedEvent.Payload[key]; !ok {
		return false
	}
	return true
}

func (m *MaskedEvent) getStringField(key string) string {
	if _, ok := m.MaskedEvent.Payload[key]; !ok {
		return ""
	}
	return m.MaskedEvent.Payload[key].(string)
}

func (m *MaskedEvent) setStringField(key string, value string) {
	m.MaskedEvent.Payload[key] = value
}

func (m *MaskedEvent) applyOperation(operation common.MaskOperation) bool {
	fmt.Println("Type:", operation.Type)
	switch operation.Type {
	case "string":
		fmt.Println("Applying string operation:", operation)
		return m.applyStringOperation(operation)
	}
	return false
}

func (m *MaskedEvent) applyStringOperation(operation common.MaskOperation) bool {
	if !m.confirmStringField(operation.Key) {
		return false
	}
	m_value := m.getStringField(operation.Key)
	switch operation.Operator {
	case "REPLACE":
		fmt.Println("Replacing", m_value, "with", operation.Input["value"].(string))
		res := operations.Replace(m_value, operation.Input["value"].(string))
		m.setStringField(operation.Key, res)
		return true
	case "RANDOM_INT":
		fmt.Println("Random int", m_value, operation.Input["upper_limit"].(string), operation.Input["lower_limit"].(string))
		res := operations.RandomInt(m_value, operation.Input["upper_limit"].(string), operation.Input["lower_limit"].(string))
		m.setStringField(operation.Key, res)
		return true
	}
	return false
}

func extractSchemaName(filePath string) string {
	// Split the path by '/'
	parts := strings.Split(filePath, "/")

	// Get the last part (file name with extension)
	fileName := parts[len(parts)-1]

	// Remove the file extension
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func loadMasks(maskMap *map[string]common.MaskConfig) {
	// Get all file paths in the schemas directory
	fmt.Println("Loading masks")
	files, err := filepath.Glob("schemas/*_mask*")
	if err != nil {
		fmt.Println("Error getting file paths:", err)
		return
	}
	fmt.Println("Files:", files)
	for _, file := range files {
		cfg, err := config.LoadMaskConfig(file)
		if err != nil {
			fmt.Println("Error loading mask config:", err)
			continue
		}
		schema := extractSchemaName(file)
		(*maskMap)[schema] = cfg
	}
	base := common.GetConfig()
	base.Masks = *maskMap
	common.SetConfig(base)
	fmt.Println("Loaded masks:", base.Masks)
}

func Synthesize(event common.Event, maskMap map[string]common.MaskConfig) {
	fmt.Printf("Synthesizing event %d times.\n", maskMap[event.Type+"_mask"].SynthAmount)
	mapObject := maskMap[event.Type+"_mask"]
	synth_amount := mapObject.SynthAmount
	for i := 0; i < synth_amount; i++ {
		fmt.Println("---------- NEXT SYNTH ----------")
		maskedEvent := new(MaskedEvent)
		maskedEvent.RawEvent = event

		// Create a deep copy of the event
		newEvent := common.Event{
			Type:    "synthed_" + string(event.Type),
			Payload: make(common.M),
		}
		for k, v := range event.Payload {
			newEvent.Payload[k] = v
		}

		maskedEvent.MaskedEvent = newEvent

		for _, operation := range mapObject.Operations {
			maskedEvent.applyOperation(operation)
		}
		// Add a unique identifier
		fmt.Println("Synthesized event:", maskedEvent.MaskedEvent)
		fmt.Println("---------- Done SYNTH ----------")
		common.MaskedStorageChan <- maskedEvent.MaskedEvent
	}
}

func createMaskedEvent(event common.Event, maskMap map[string]common.MaskConfig) MaskedEvent {
	mapObject := maskMap[event.Type+"_mask"]
	maskedEvent := MaskedEvent{
		RawEvent:    event,
		MaskedEvent: event,
		Type:        event.Type,
	}
	maskedEvent.MaskedEvent.Type = "masked_" + event.Type
	for _, operation := range mapObject.Operations {
		maskedEvent.applyOperation(operation)
	}

	fmt.Println("About to synth:", mapObject.SynthAmount)
	if mapObject.SynthAmount > 0 {
		go Synthesize(maskedEvent.RawEvent, maskMap)
	}

	return maskedEvent
}

func Mask(maskChan *chan common.Event, mask_storage bool) {
	fmt.Println("Masking service started")
	loadMasks(&MaskMap)
	base := common.GetConfig()
	masks := base.Masks
	fmt.Println("Masks:", masks)
	for event := range *maskChan {
		maskedEvent := createMaskedEvent(event, masks)
		if mask_storage {
			common.MaskedStorageChan <- maskedEvent.MaskedEvent
		}
	}
}

// Refactor this to just use common.Event as output
func MaskSingleEvent(event common.Event, maskMap map[string]common.MaskConfig) MaskedEvent {
	fmt.Println("Masking single event:", event)
	masks := common.GetConfig().Masks
	maskedEvent := createMaskedEvent(event, masks)
	return maskedEvent
}
