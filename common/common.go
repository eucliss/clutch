package common // or your main package name

import (
	"clutch/services/model"
	"fmt"
	"sync"
)

var (
	// EventChan is a channel for sending events throughout the program
	EventChan         = make(chan Event, 1000)
	Pipeline          = make(chan Event, 1000)
	MaskChan          = make(chan Event, 1000)
	StorageChan       = make(chan Event, 1000)
	MaskedStorageChan = make(chan Event, 1000)
	SynthChan         = make(chan Event, 1000)
	ChatChan          = make(chan Event, 1000)

	// ErrorChan is a channel for sending errors throughout the program
	ErrorChan = make(chan error, 100)

	// GlobalConfig holds the program-wide configuration
	GlobalConfig Config

	// Mutex for safe concurrent access to shared resources
	globalMutex sync.Mutex
)

type M map[string]interface{}

type Store interface {
	InsertDocument(index string, body map[string]interface{})
	Query(index string, query string) (r map[string]interface{})
	DeleteIndex(index string)
	Initialize()
	GetResults(searchResult map[string]interface{}) (res []map[string]interface{})
}

// Event struct for your event channel
type Event struct {
	Type    string
	Payload M
}

type MaskConfig struct {
	SynthAmount int             `yaml:"synthetic_count"`
	Operations  []MaskOperation `yaml:"masks"`
}

type MaskOperation struct {
	Key      string `yaml:"key"`
	Operator string `yaml:"operator"`
	Input    M      `yaml:"input"`
	Type     string `yaml:"type"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

// Struct to represent the Database section
type DatabaseConfig struct {
	Type         string `yaml:"type"`
	CertLocation string `yaml:"cert_location"`
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	NewIndex     string `yaml:"new_index_on_launch"`
}

// Struct to represent the full configuration
type Config struct {
	Server   ServerConfig          `yaml:"server"`
	Database DatabaseConfig        `yaml:"database"`
	Services []string              `yaml:"services"`
	Store    Store                 `yaml:"store"`
	Masks    map[string]MaskConfig `yaml:"masks"`
	Model    model.Model           `yaml:"model"`
}

func GetConfigAddress() *Config {
	globalMutex.Lock()
	defer globalMutex.Unlock()
	return &GlobalConfig
}

// GetConfig returns a copy of the global configuration
func GetConfig() Config {
	globalMutex.Lock()
	defer globalMutex.Unlock()
	return GlobalConfig
}

func (c *Config) SetDbConfig(cfg DatabaseConfig) {
	globalMutex.Lock()
	defer globalMutex.Unlock()
	c.Database = cfg
}

func (c *Config) SetStoreConfig(cfg Store) {
	globalMutex.Lock()
	defer globalMutex.Unlock()
	c.Store = cfg
}

func (c *Config) SetModelConfig(cfg model.Model) {
	globalMutex.Lock()
	defer globalMutex.Unlock()
	c.Model = cfg
}

// SetConfig updates the global configuration
func SetConfig(cfg Config) {
	globalMutex.Lock()
	defer globalMutex.Unlock()
	GlobalConfig = cfg
}

func FlattenMap(result map[string]interface{}, prefix string, m map[string]interface{}) {
	for k, v := range m {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}

		switch val := v.(type) {
		case map[string]interface{}:
			FlattenMap(result, key, val)
		case []interface{}:
			for i, arrayVal := range val {
				if subMap, ok := arrayVal.(map[string]interface{}); ok {
					FlattenMap(result, fmt.Sprintf("%s.%d", key, i), subMap)
				} else {
					result[fmt.Sprintf("%s.%d", key, i)] = arrayVal
				}
			}
		default:
			result[key] = v
		}
	}
}
