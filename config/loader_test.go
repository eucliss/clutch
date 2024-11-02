package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary config file for testing
	tempFile, err := os.CreateTemp("", "config.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write test configuration to the temp file
	testConfig := []byte(`
server:
  host: "localhost"
  port: "8080"
database:
  type: "elastic"
  cert_location: "db/http_ca.crt"
  host: "localhost"
  port: "9200"
  user: "elastic"
  password: "testpassword"
  new_index_on_launch: "test-index"
`)
	if _, err := tempFile.Write(testConfig); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	// Temporarily replace the config file path
	oldConfigPath := ConfigPath // Assume ConfigPath is exported from the loader package
	ConfigPath = tempFile.Name()
	defer func() { ConfigPath = oldConfigPath }()

	// Test the LoadConfig function
	config, err := LoadCommonConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	// Add assertions to check if the config was loaded correctly
	if config.Server.Host != "localhost" {
		t.Errorf("Expected server host to be 'localhost', got %s", config.Server.Host)
	}
	// Add more assertions as needed
}

// TestLoadConfigFileNotFound tests the error case when the config file is not found
func TestLoadConfigFileNotFound(t *testing.T) {
	// Ensure the file doesn't exist
	os.Remove("config.yaml")

	_, err := LoadCommonConfig()
	if err == nil {
		t.Error("LoadConfig() expected an error, got nil")
	}
}
