package mageutil

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"runtime"
)

var (
	serviceBinaries    map[string]int
	toolBinaries       []string
	MaxFileDescriptors int
)

type Config struct {
	ServiceBinaries    map[string]int `yaml:"serviceBinaries"`
	ToolBinaries       []string       `yaml:"toolBinaries"`
	MaxFileDescriptors int            `yaml:"maxFileDescriptors"`
}

func InitForSSC() {
	yamlFile, err := ioutil.ReadFile("start-config.yml")
	if err != nil {
		log.Fatalf("error reading YAML file: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("error unmarshalling YAML: %v", err)
	}

	adjustedBinaries := make(map[string]int)
	for binary, count := range config.ServiceBinaries {
		if runtime.GOOS == "windows" {
			binary += ".exe"
		}
		adjustedBinaries[binary] = count
	}

	serviceBinaries = adjustedBinaries
	toolBinaries = config.ToolBinaries
	MaxFileDescriptors = config.MaxFileDescriptors
}
