package main

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/manifoldco/promptui"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"text/template"
)

const (
	// configPath is a variable pointing to the jumpr config.yaml file from the home directory for your user
	configPath = "/.jumpr/config.yaml"
)

var (
	data []cluster
	//go:embed "templates/*"
	templates embed.FS
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

func main() {
	switch len(os.Args) {
	case 1:
		promptUser()
	case 2:
		inlineLogin(os.Args[1])
	default:
		log.Fatal("Unsupported number of arguments. Use either 1 or none.")
	}
}

func promptUser() {
	searcher := func(input string, index int) bool {
		return fuzzy.Match(input, data[index].ClusterNickname) || fuzzy.Match(input, data[index].ClusterName)
	}

	promptTemplate := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "{{ .ClusterNickname | blue }} ({{ .ClusterName | faint }})",
		Inactive: "{{ .ClusterNickname | faint }}",
		Selected: "Selected: {{ .ClusterNickname | blue }}",
		Details: `
--------- Cluster ----------
{{ "Nickname:" | faint }}	 {{ .ClusterNickname }}
{{ "Name:" | faint }}	 {{ .ClusterName }}
{{ "Region:" | faint }}	 {{ .AWSRegion }}
{{ "Profile:" | faint }}	 {{ .AWSProfile }}`,
	}

	prompt := promptui.Select{
		Label:     "Select a cluster",
		Items:     data,
		Size:      8,
		Templates: promptTemplate,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()
	if err != nil {
		log.Fatal(err)
	}

	data[i].ssoLogin()
	data[i].getCert()
	data[i].getEndpoint()
	data[i].generateKubeconfig()
	data[i].printExports()
}

func inlineLogin(input string) {
	for _, cluster := range data {
		if cluster.ClusterNickname == input || cluster.ClusterName == input {
			cluster.ssoLogin()
			cluster.getCert()
			cluster.getEndpoint()
			cluster.generateKubeconfig()
			cluster.printExports()
			return
		}
	}
	log.Fatalf("%v was not recognized as any cluster in %v", input, os.Getenv("HOME")+configPath)
}

type cluster struct {
	ClusterName     string `yaml:"clusterName"`
	ClusterNickname string `yaml:"clusterNickname"`
	AWSProfile      string `yaml:"awsProfile"`
	AWSRegion       string `yaml:"awsRegion"`
	AWSAccountId    string `yaml:"awsAccountId"`
	CertificateData string // not provided by config.yaml
	ClusterEndpoint string // not provided by config.yaml
	KubeconfigPath  string // not provided by config.yaml
}

func (c *cluster) ssoLogin() {
	cmd := exec.Command("aws", "sso", "login", "--profile", c.AWSProfile)

	err := cmd.Run()
	if err != nil {
		log.Fatalf("Failed to login into profile %v, %v", c.AWSProfile, err)
	}

	if err = os.Setenv("AWS_DEFAULT_PROFILE", c.AWSProfile); err != nil {
		log.Fatalf("Failed to set AWS_DEFAULT_PROFILE for further operations, %v", err)
	}
}

func (c *cluster) generateKubeconfig() {
	tmpl, err := template.New("kubeconfig.tpl").ParseFS(templates, "templates/kubeconfig.tpl")
	if err != nil {
		log.Fatalf("Failed to parse kubeconfig template, %v", err)
	}

	b := new(bytes.Buffer)

	err = tmpl.Execute(b, c)
	if err != nil {
		log.Fatalf("Failed to execute kubeconfig template, %v", err)
	}

	kubePath := "/tmp/kubeconfig-" + c.ClusterName

	c.KubeconfigPath = kubePath

	f, err := os.Create(kubePath)
	if err != nil {
		log.Fatalf("Failed creating a file at %v, %v", kubePath, err)
	}
	defer f.Close()

	_, err = f.WriteString(b.String())
}

func (c *cluster) getEndpoint() {
	output, err := exec.Command("aws", "eks", "describe-cluster", "--name", c.ClusterName, "--query", "cluster.endpoint", "--output", "text", "--region", c.AWSRegion).Output()
	if err != nil {
		log.Fatalf("Failed to get cluster endpoint, %v", err)
	}
	c.ClusterEndpoint = string(output)
}

func (c *cluster) getCert() {
	output, err := exec.Command("aws", "eks", "describe-cluster", "--name", c.ClusterName, "--query", "cluster.certificateAuthority.data", "--output", "text", "--region", c.AWSRegion).Output()
	if err != nil {
		log.Fatalf("Failed to get cluster certificate data, %v", err)
	}
	c.CertificateData = string(output)
}

func (c *cluster) printExports() {
	fmt.Printf(`
Paste the commands below into your CLI to remain in this profile/region:
------------------------------------------------------------------------
export AWS_DEFAULT_PROFILE="%v"
export AWS_REGION="%v"
export KUBECONFIG="%v"
export http_proxy=http://10.122.32.110:3128
`, c.AWSProfile, c.AWSRegion, c.KubeconfigPath)
}
