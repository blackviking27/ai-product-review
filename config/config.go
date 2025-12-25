package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	Server   ServerConfig   `json:"server"`
	AI       AIConfig       `json:"ai"`
	Scrapper ScrapperConfig `json:"scrappers"`
}

type ServerConfig struct {
	Port    string `json:"port"`
	Timeout int    `json:"timeout_seconds"`
}

type AIConfig struct {
	Provider string `json:"provider"`
	Model    string `json:"model"`
	ApiKey   string `json:"api_key"`
}

type ScrapperConfig struct {
	UserAgent string         `json:"user_agent"`
	Amazon    PlatformConfig `json:"amazon"`
	Flipkart  PlatformConfig `json:"flipkart"`
}

type PlatformConfig struct {
	ReviewContainer string     `json:"review_container"`
	ReviewText      string     `json:"review_text"`
	ReviewRating    string     `json:"review_rating"`
	ReviewLink      ReviewLink `json:"review_link"`
}

type ReviewLink struct {
	Host     string `json:"host"`
	Path     string `json:"path"`
	FullLink string `json:"full_link"`
}

func LoadConfig() (*Config, error) {

	path := filepath.Join("config", "config.json")

	fileBytes, err := os.ReadFile(path)

	if err != nil {
		log.Fatal("Unable to load config.json file")
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(fileBytes, &config); err != nil {
		log.Fatal("Unable to load config values")
		return nil, err
	}

	return &config, nil
}
