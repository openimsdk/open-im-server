package config

import (
	"flag"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Directory string   `yaml:"directory"`
	FileTypes []string `yaml:"file_types"`
	Languages []string `yaml:"languages"`
}

func ParseConfig() (Config, error) {
	var configPath string
	flag.StringVar(&configPath, "config", "./", "Path to config file")
	flag.Parse()

	var config Config
	if configPath != "" {
		configFile, err := os.ReadFile(configPath)
		if err != nil {
			return Config{}, err
		}
		if err := yaml.Unmarshal(configFile, &config); err != nil {
			return Config{}, err
		}
	} else {
		log.Fatal("Config file must be provided")
	}
	return config, nil
}
