package cli

import (
	"log"
	"os"
)

func InlineLogin(input string) {
	for _, cluster := range data {
		if cluster.ClusterNickname == input || cluster.ClusterName == input {
			cluster.SsoLogin()
			cluster.GetCert()
			cluster.GetEndpoint()
			cluster.GenerateKubeconfig()
			cluster.PrintExports()
			return
		}
	}
	log.Fatalf("%v was not recognized as any cluster in %v", input, os.Getenv("HOME")+configPath)
}
