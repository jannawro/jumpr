package internal

import (
	"log"
	"os"
)

func InlineLogin(input string) {
	for _, cluster := range jumprConfig.Clusters {
		if cluster.Nickname == input || cluster.Name == input {
			cluster.SsoLogin()
			cluster.GetClusterInfo()
			cluster.GenerateKubeconfig()
			cluster.PrintExports()
			return
		}
	}
	log.Fatalf("%v was not recognized as any cluster in %v", input, os.Getenv("HOME")+configPath)
}
