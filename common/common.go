package common // or your main package name

import (
	"clutch/store"
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

	// ErrorChan is a channel for sending errors throughout the program
	ErrorChan = make(chan error, 100)

	// GlobalConfig holds the program-wide configuration
	GlobalConfig Config

	// Mutex for safe concurrent access to shared resources
	globalMutex sync.Mutex
)

type M map[string]interface{}

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

type ModelConfig struct {
	URL       string `yaml:"url"`
	ModelName string `yaml:"model_name"`
}

// Struct to represent the full configuration
type Config struct {
	Server   ServerConfig          `yaml:"server"`
	Database DatabaseConfig        `yaml:"database"`
	Services []string              `yaml:"services"`
	Store    store.Store           `yaml:"store"`
	Masks    map[string]MaskConfig `yaml:"masks"`
	Model    ModelConfig           `yaml:"model"`
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

func (c *Config) SetStoreConfig(cfg store.Store) {
	globalMutex.Lock()
	defer globalMutex.Unlock()
	c.Store = cfg
}

// SetConfig updates the global configuration
func SetConfig(cfg Config) {
	globalMutex.Lock()
	defer globalMutex.Unlock()
	GlobalConfig = cfg
}
