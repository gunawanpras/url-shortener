package config

import (
	"log"
	"os"

	"github.com/gunawanpras/url-shortener/cache"
	"gopkg.in/yaml.v2"
)

type (
	Server struct {
		Port string `yaml:"port"`
	}

	URLShortener struct {
		BaseURL string `yaml:"base_url"`
	}

	Config struct {
		Server       Server            `yaml:"server"`
		URLShortener URLShortener      `yaml:"url_shortener"`
		Cache        cache.CacheConfig `yaml:"cache"`
	}
)

func LoadConfig(filePath string) Config {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Error opening config file:", err)
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		log.Fatal("Error decoding config file:", err)
	}

	return config
}
