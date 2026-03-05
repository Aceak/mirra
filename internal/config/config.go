package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config configuration structure
type Config struct {
	Server     ServerConfig     `json:"server"`
	Share      ShareConfig      `json:"share"`
	Appearance AppearanceConfig `json:"appearance"`
}

// ServerConfig server configuration
type ServerConfig struct {
	Name    string `json:"name"`
	Host    string `json:"host"`
	Port    string `json:"port"`
	Favicon string `json:"favicon"`
}

// ShareConfig share configuration
type ShareConfig struct {
	RootPath string `json:"root_path"`
}

// AppearanceConfig appearance configuration
type AppearanceConfig struct {
	Theme      string `json:"theme"`
	ShowHidden bool   `json:"show_hidden"`
}

var config Config

// LoadConfig loads configuration file and returns Config instance
// configPath: path to config file, defaults to config.json
func LoadConfig(configPath string) (*Config, error) {
	// If config file does not exist, create default configuration
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Println("config.json not found, creating default configuration...")

		// Get current working directory as default root path
		cwd, err := os.Getwd()
		if err != nil {
			cwd = "."
		}

		defaultConfig := Config{
			Server: ServerConfig{
				Name:    "Mirra",
				Host:    "0.0.0.0",
				Port:    "8080",
				Favicon: "",
			},
			Share: ShareConfig{
				RootPath: cwd,
			},
			Appearance: AppearanceConfig{
				Theme:      "auto",
				ShowHidden: false,
			},
		}

		// Write default configuration
		file, err := os.Create(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to create config.json: %v", err)
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		if encodeErr := encoder.Encode(defaultConfig); encodeErr != nil {
			return nil, fmt.Errorf("failed to write default config: %v", encodeErr)
		}

		fmt.Printf("Default config.json created with root path: %s\n", cwd)

		// Reopen file for reading
		file, err = os.Open(configPath)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		if decodeErr := decoder.Decode(&config); decodeErr != nil {
			return nil, decodeErr
		}
		return &config, nil
	}

	// Config file exists, load normally
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
