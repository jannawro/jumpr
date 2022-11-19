package cli

import (
	cl "github.com/jannawro/jumpr/internal/clusterLogin"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
)

const (
	// configPath is a variable pointing to the jumpr config.yaml file from the home directory for your user
	configPath = "/.jumpr/config.yaml"
)

var (
	data []cl.Cluster
)

func init() {
	configBody, err := ioutil.ReadFile(os.Getenv("HOME") + configPath)
	if err != nil {
		log.Fatalf("Failed to read from your %v file, %v", os.Getenv("HOME")+configPath, err)
	}

	err = yaml.Unmarshal(configBody, &data)
	if err != nil {
		log.Fatalf("Failed to extract data form your %v file, %v", os.Getenv("HOME")+configPath, err)
	}
}
