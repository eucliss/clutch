package common // or your main package name

import (
	"fmt"
	"sync"
)

var (
	// Channels for sending events throughtout the services
	// Base channel for ingest
	EventChan = make(chan Event, 1000)
	// All events go through the pipeline
	Pipeline = make(chan Event, 1000)
	// Masking Channel for masking events
	MaskChan = make(chan Event, 1000)
	// Storage Channel for storing events
	StorageChan = make(chan Event, 1000)
	// Masked Storage Channel for storing masked events
	MaskedStorageChan = make(chan Event, 1000)
	// Synthesis Channel for synthesis events
	SynthChan = make(chan Event, 1000)
	// Chat Channel for chat events
	ChatChan = make(chan Event, 1000)

	// ErrorChan is a channel for sending errors throughout the program
	ErrorChan = make(chan error, 100)

	// GlobalConfig holds the program-wide configuration
	GlobalConfig Config

	// Mutex for safe concurrent access to shared resources
	globalMutex sync.Mutex
)

// JSON Map type
type M map[string]interface{}

type ModelInterface interface {
	GenerateEmbeddings(text string) ([][]float32, error)
	QueryWithContext(query string, ctx string) (string, error)
	Start()
	GetModelName() string
}

type Store interface {
	InsertDocument(index string, body map[string]interface{})
	Query(index string, query string) (r map[string]interface{})
	DeleteIndex(index string)
	Initialize()
	GetResults(searchResult map[string]interface{}) (res []map[string]interface{})
}

type Masker interface {
	Mask(event Event, maskConfig MaskConfig) Event
	Synthesize(event Event, maskConfig MaskConfig) Event
	// Refactor this to just use common.Event as output
	// MaskSingleEvent(event Event, maskConfig MaskConfig) MaskedEvent
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

type BaseModelConfig struct {
	URL               string `yaml:"url"`
	EmbedderURL       string `yaml:"embedder_url"`
	EmbedderModelName string `yaml:"embedder_model_name"`
	ModelName         string `yaml:"model_name"`
	BasePrompt        string `yaml:"base_prompt"`
}

// Struct to represent the full configuration
type Config struct {
	Server      ServerConfig          `yaml:"server"`
	Database    DatabaseConfig        `yaml:"database"`
	Services    []string              `yaml:"services"`
	Masks       map[string]MaskConfig `yaml:"masks"`
	ModelConfig BaseModelConfig       `yaml:"model"`
	Model       ModelInterface        `yaml:"-"`
	Store       Store                 `yaml:"-"`
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

// SetConfig updates the global configuration
func SetConfig(cfg Config) {
	globalMutex.Lock()
	defer globalMutex.Unlock()
	GlobalConfig = cfg
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

func (c *Config) SetModelConfig(cfg ModelInterface) {
	globalMutex.Lock()
	defer globalMutex.Unlock()
	c.Model = cfg
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
