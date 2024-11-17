package common

import (
	"os"
	"reflect"
	"testing"
)

// Mock implementations for interfaces
type MockModelInterface struct{}

func (m *MockModelInterface) GenerateEmbeddings(text string) ([][]float32, error) { return nil, nil }
func (m *MockModelInterface) QueryWithContext(query string, ctx string) (string, error) {
	return "NICE!", nil
}
func (m *MockModelInterface) Start()               {}
func (m *MockModelInterface) GetModelName() string { return "mock_model" }

func TestModelInterface(t *testing.T) {
	var model ModelInterface
	model = &MockModelInterface{}
	if model.GetModelName() != "mock_model" {
		t.Errorf("GetModelName() = %v, want %v", model.GetModelName(), "mock_model")
	}
}

type MockStore struct{}

func (s *MockStore) InsertDocument(index string, body map[string]interface{}) {}
func (s *MockStore) Query(index string, query string) (r map[string]interface{}) {
	return map[string]interface{}{
		"test": "test",
	}
}
func (s *MockStore) DeleteIndex(index string) {}
func (s *MockStore) Initialize()              {}
func (s *MockStore) GetResults(searchResult map[string]interface{}) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"test": "test",
		},
	}
}

func TestStoreInterface(t *testing.T) {
	var store Store
	store = &MockStore{}
	results := store.GetResults(map[string]interface{}{})
	if results[0]["test"] != "test" {
		t.Errorf("GetResults() = %v, want %v", results[0]["test"], "test")
	}
}

var (
	testConfig Config
)

func TestMain(m *testing.M) {
	// Setup
	testConfig = Config{
		Server: ServerConfig{
			Host: "localhost",
			Port: "8080",
		},
		Database: DatabaseConfig{
			Type: "postgres",
			Host: "localhost",
			Port: "5432",
		},
		Services: []string{"service1", "service2"},
	}
	// Set initial global config
	SetConfig(testConfig)

	// Run tests
	code := m.Run()

	// Exit with test result code
	os.Exit(code)
}

func TestGetConfigAddress(t *testing.T) {
	configAddr := GetConfigAddress()
	if configAddr == nil {
		t.Error("GetConfigAddress() returned nil")
	}
	if !reflect.DeepEqual(*configAddr, testConfig) {
		t.Errorf("GetConfigAddress() = %v, want %v", *configAddr, testConfig)
	}
}

func TestGetConfig(t *testing.T) {
	config := GetConfig()
	if !reflect.DeepEqual(config, testConfig) {
		t.Errorf("GetConfig() = %v, want %v", config, testConfig)
	}
}

func TestSetConfig(t *testing.T) {
	oldConfig := testConfig
	newConfig := Config{}
	SetConfig(newConfig)
	config := GetConfig()

	if !reflect.DeepEqual(config, newConfig) {
		t.Errorf("Config not set: %v, want %v", config, newConfig)
	}
	if reflect.DeepEqual(oldConfig, config) {
		t.Errorf("Old config still set: %v, want %v", oldConfig, config)
	}
}

func TestSetDbConfig(t *testing.T) {
	oldDbConfig := testConfig.Database
	if oldDbConfig.Host != "localhost" {
		t.Error("Old database config not set")
	}
	newDbConfig := DatabaseConfig{
		Host: "new_host",
	}
	testConfig.SetDbConfig(newDbConfig)

	if !reflect.DeepEqual(testConfig.Database, newDbConfig) {
		t.Errorf("Database config not set: %v, want %v", testConfig.Database, newDbConfig)
	}
	if testConfig.Database.Host != "new_host" {
		t.Errorf("Database host not updated: %v, want %v", testConfig.Database.Host, "new_host")
	}
}

func TestSetStoreConfig(t *testing.T) {
	newStore := &MockStore{}
	testConfig.SetStoreConfig(newStore)

	store := testConfig.Store
	res := store.Query("foo", "bar")
	if res["test"] != "test" {
		t.Errorf("Store not set: %v, want %v", res["test"], "test")
	}
}

func TestSetModelConfig(t *testing.T) {
	newModel := &MockModelInterface{}
	testConfig.SetModelConfig(newModel)

	model := testConfig.Model
	res, _ := model.QueryWithContext("foo", "bar")
	if res != "NICE!" {
		t.Errorf("Model not set: %v, want %v", res, "NICE!")
	}
}

func TestFlattenMap(t *testing.T) {
	input := map[string]interface{}{
		"foo": "bar",
		"baz": map[string]interface{}{
			"test": "test",
		},
	}
	expected := map[string]interface{}{
		"foo":      "bar",
		"baz.test": "test",
	}
	result := map[string]interface{}{}
	FlattenMap(result, "", input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("FlattenMap() = %v, want %v", result, expected)
	}
}

func TestFlattenMapPrefix(t *testing.T) {
	input := map[string]interface{}{
		"foo": "bar",
		"baz": map[string]interface{}{
			"test": "test",
		},
	}
	expected := map[string]interface{}{
		"prefix.foo":      "bar",
		"prefix.baz.test": "test",
	}
	result := map[string]interface{}{}
	FlattenMap(result, "prefix", input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("FlattenMap() = %v, want %v", result, expected)
	}
}
