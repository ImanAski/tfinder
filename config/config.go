package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Dir     string            `json:"dir"`
	Ignore  []string          `json:"ignore"`
	Pattern []string          `json:"pattern"`
	Colors  map[string]string `json:"colors"`
}

var defaultConfig = Config{
	Dir:     ".",
	Ignore:  []string{".git", "node_modules", "vendor"},
	Pattern: []string{"TODO", "FIXME", "DO", "BUG"},
	Colors: map[string]string{
		"TODO":  "yellow",
		"FIXME": "red",
		"DO":    "green",
		"BUG":   "red",
	},
}

func LoadConfig() Config {
	configFile := "tfinder.json"

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return defaultConfig
	}

	file, err := os.Open(configFile)
	if err != nil {
		fmt.Printf("Warning: could not open config file: %v. using defaults\n", err)
		return defaultConfig
	}

	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		fmt.Printf("Invalid config, using defaults: %v\n", err)
		return defaultConfig
	}

	return config
}

func CreateConfig() error {
	configFile := "tfinder.json"
	file, err := os.Create(configFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(defaultConfig); err != nil {
		return err
	}

	return nil
}
