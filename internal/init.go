package internal

import (
	"os"

	"gopkg.in/yaml.v3"
)

const (
	// configPath is a variable pointing to the jumpr config.yaml file from the home directory for your user
	configPath = "/.jumpr/config.yaml"
)

type JumprConfig struct {
	Clusters []Cluster
}

var (
	jumprConfig JumprConfig
)

func init() {
	configBody, err := os.ReadFile(os.Getenv("HOME") + configPath)
	check("Failed at reading form your ~/.jumpr/config.yaml file:", err)

	err = yaml.Unmarshal(configBody, &jumprConfig)
	check("Failed at extracting data from your ~/.jumpr/config.yaml file:", err)
}
