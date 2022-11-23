package clusterLogin

import (
	"bytes"
	"embed"
	"log"
	"os"
	"os/exec"
	"text/template"
)

var (
	//go:embed "templates/*"
	templates embed.FS
)

type Cluster struct {
	Name            string `yaml:"clusterName"`
	Nickname        string `yaml:"clusterNickname"`
	Profile         string `yaml:"awsProfile"`
	Region          string `yaml:"awsRegion"`
	AccountId       string `yaml:"awsAccountId"`
	Proxy           string `yaml:"proxy"`
	CertificateData string // not provided by config.yaml
	ClusterEndpoint string // not provided by config.yaml
	KubeconfigPath  string // not provided by config.yaml
}

func (c *Cluster) SsoLogin() {
	cmd := exec.Command("aws", "sso", "login", "--profile", c.Profile)

	err := cmd.Run()
	if err != nil {
		log.Fatalf("Failed to login into profile %v, %v", c.Profile, err)
	}

	if err = os.Setenv("AWS_DEFAULT_PROFILE", c.Profile); err != nil {
		log.Fatalf("Failed to set AWS_DEFAULT_PROFILE for further operations, %v", err)
	}
}

func (c *Cluster) GenerateKubeconfig() {
	tmpl, err := template.New("kubeconfig.gotmpl").ParseFS(templates, "templates/kubeconfig.gotmpl")
	if err != nil {
		log.Fatalf("Failed to parse kubeconfig template, %v", err)
	}

	b := new(bytes.Buffer)

	err = tmpl.Execute(b, c)
	if err != nil {
		log.Fatalf("Failed to execute kubeconfig template, %v", err)
	}

	kubePath := "/tmp/kubeconfig-" + c.Name

	c.KubeconfigPath = kubePath

	f, err := os.Create(kubePath)
	if err != nil {
		log.Fatalf("Failed creating a file at %v, %v", kubePath, err)
	}
	defer f.Close()

	_, err = f.WriteString(b.String())
}

func (c *Cluster) GetEndpoint() {
	output, err := exec.Command("aws", "eks", "describe-cluster", "--name", c.Name,
		"--query", "cluster.endpoint", "--output", "text", "--region", c.Region).Output()
	if err != nil {
		log.Fatalf("Failed to get Cluster endpoint, %v", err)
	}
	c.ClusterEndpoint = string(output)
}

func (c *Cluster) GetCert() {
	output, err := exec.Command("aws", "eks", "describe-cluster", "--name", c.Name,
		"--query", "cluster.certificateAuthority.data", "--output", "text", "--region", c.Region).Output()
	if err != nil {
		log.Fatalf("Failed to get cluster certificate data, %v", err)
	}
	c.CertificateData = string(output)
}

func (c *Cluster) PrintExports() {
	tmpl, err := template.New("exports.gotmpl").ParseFS(templates, "templates/exports.gotmpl")
	if err != nil {
		log.Fatalf("Failed to parse kubeconfig template, %v", err)
	}

	err = tmpl.Execute(os.Stdout, c)
	if err != nil {
		log.Fatalf("Failed to execute exports template, %v", err)
	}
}
