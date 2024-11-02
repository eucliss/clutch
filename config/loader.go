package config

import (
	"clutch/common"
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v3"
)

var ConfigPath string

func LoadMaskConfig(path string) (common.MaskConfig, error) {
	fmt.Println("Loading mask config:", path)
	file, err := os.Open(path)
	if err != nil {
		return common.MaskConfig{}, fmt.Errorf("error reading YAML file: %w", err)
	}
	defer file.Close()

	var config common.MaskConfig

	// Unmarshal the YAML data into the Config struct
	decoder := yaml.NewDecoder(file)
	if err = decoder.Decode(&config); err != nil {
		return common.MaskConfig{}, err
	}
	return config, err
}

func LoadCommonConfig() (common.Config, error) {
	fmt.Println("Loading base config")
	if ConfigPath == "" {
		ConfigPath = "config.yaml"
	}
	file, err := os.Open(ConfigPath)
	if err != nil {
		return common.Config{}, fmt.Errorf("error reading YAML file: %w", err)
	}
	defer file.Close()

	var config common.Config

	// Unmarshal the YAML data into the Config struct
	decoder := yaml.NewDecoder(file)
	if err = decoder.Decode(&config); err != nil {
		return common.Config{}, err
	}
	return config, err
}
