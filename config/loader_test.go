package config

import (
	"clutch/common"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Create a temporary config file for testing
	tempFile, err := os.CreateTemp("", "config.yaml")
	if err != nil {
		panic(err)
	}
	defer os.Remove(tempFile.Name())

	// Write test configuration to the temp file
	testConfig := []byte(`
server:
  host: "localhost"
  port: "8080"
model:
  url: "http://localhost:11434"
  model_name: "llama3.2"
  embedder_url: "http://localhost:11434"
  embedder_model_name: "nomic-embed-text"
  base_prompt: "base prompt"
database:
  type: "qdrant"
  cert_location: "http_ca.crt"
  host: "localhost"
  port: "6334"
  user: "qdrant"
  password: "password"
services:
  - storage
  - model
  - masking
  - mask_storage
`)

	if _, err := tempFile.Write(testConfig); err != nil {
		panic(err)
	}
	tempFile.Close()

	// // Temporarily replace the config file path
	oldConfigPath := ConfigPath // Assume ConfigPath is exported from the loader package
	ConfigPath = tempFile.Name()
	defer func() { ConfigPath = oldConfigPath }()

	// Run tests
	code := m.Run()

	// Exit with test result code
	os.Exit(code)
}

func TestInitializeConfig(t *testing.T) {
	success, err := InitializeConfig()
	if err != nil {
		t.Fatalf("InitializeConfig() error = %v", err)
	}
	if !success {
		t.Error("InitializeConfig() expected true, got false")
	}
	config := common.GetConfig()
	if config.Server.Host != "localhost" {
		t.Errorf("Expected server host to be 'localhost', got %s", config.Server.Host)
	}
}

func TestLoadConfig(t *testing.T) {
	// Test the LoadConfig function
	config, err := LoadCommonConfig()
	if err != nil {
		t.Fatalf("LoadCommonConfig() error = %v", err)
	}
	// Add assertions to check if the config was loaded correctly
	if config.Server.Host != "localhost" {
		t.Errorf("Expected server host to be 'localhost', got %s", config.Server.Host)
	}
	if config.ModelConfig.ModelName != "llama3.2" {
		t.Errorf("Expected model name to be 'llama3.2', got %s", config.ModelConfig.ModelName)
	}
	if config.Database.Type != "qdrant" {
		t.Errorf("Expected database type to be 'qdrant', got %s", config.Database.Type)
	}
}

func TestLoadMaskConfig(t *testing.T) {
	config, err := LoadMaskConfig("../schemas/testing.yaml")
	if err != nil {
		t.Fatalf("LoadMaskConfig() error = %v", err)
	}
	if len(config.Operations) != 2 {
		t.Errorf("Expected 2 masks, got %d", len(config.Operations))
	}
	if config.SynthAmount != 1 {
		t.Errorf("Expected 1 synthetic count, got %d", config.SynthAmount)
	}
}
