package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config 配置结构
type Config struct {
	Server     ServerConfig     `json:"server"`
	Share      ShareConfig      `json:"share"`
	Appearance AppearanceConfig `json:"appearance"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Name    string `json:"name"`
	Host    string `json:"host"`
	Port    string `json:"port"`
	Favicon string `json:"favicon"`
}

// ShareConfig 共享配置
type ShareConfig struct {
	RootPath string `json:"root_path"`
}

// AppearanceConfig 外观配置
type AppearanceConfig struct {
	Theme      string `json:"theme"`
	ShowHidden bool   `json:"show_hidden"`
}

var config Config

// LoadConfig 加载配置文件并返回 Config 实例
func LoadConfig() (*Config, error) {
	// 如果配置文件不存在，创建默认配置
	if _, err := os.Stat("config.json"); os.IsNotExist(err) {
		fmt.Println("config.json not found, creating default configuration...")

		// 获取当前工作目录作为默认根目录
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

		// 写入默认配置
		file, err := os.Create("config.json")
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

		// 重新打开文件进行读取
		file, err = os.Open("config.json")
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

	// 配置文件存在，正常加载
	file, err := os.Open("config.json")
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
