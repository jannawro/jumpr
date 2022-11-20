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
	name            string `yaml:"name"`
	nickname        string `yaml:"clusterNickname"`
	profile         string `yaml:"awsProfile"`
	region          string `yaml:"awsRegion"`
	accountId       string `yaml:"awsAccountId"`
	proxy           string `yaml:"proxy"`
	certificateData string // not provided by config.yaml
	clusterEndpoint string // not provided by config.yaml
	kubeconfigPath  string // not provided by config.yaml
}

func (c *Cluster) Name() string {
	return c.name
}

func (c *Cluster) Nickname() string {
	return c.nickname
}

func (c *Cluster) Profile() string {
	return c.profile
}

func (c *Cluster) Region() string {
	return c.region
}

func (c *Cluster) AccountId() string {
	return c.accountId
}

func (c *Cluster) Proxy() string {
	return c.proxy
}

func (c *Cluster) CertificateData() string {
	return c.certificateData
}

func (c *Cluster) ClusterEndpoint() string {
	return c.clusterEndpoint
}

func (c *Cluster) KubeconfigPath() string {
	return c.kubeconfigPath
}

func (c *Cluster) FilterValue() string {
	return c.nickname
}

func (c *Cluster) SsoLogin() {
	cmd := exec.Command("aws", "sso", "login", "--profile", c.profile)

	err := cmd.Run()
	if err != nil {
		log.Fatalf("Failed to login into profile %v, %v", c.profile, err)
	}

	if err = os.Setenv("AWS_DEFAULT_PROFILE", c.profile); err != nil {
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

	kubePath := "/tmp/kubeconfig-" + c.name

	c.kubeconfigPath = kubePath

	f, err := os.Create(kubePath)
	if err != nil {
		log.Fatalf("Failed creating a file at %v, %v", kubePath, err)
	}
	defer f.Close()

	_, err = f.WriteString(b.String())
}

func (c *Cluster) GetEndpoint() {
	output, err := exec.Command("aws", "eks", "describe-cluster", "--name", c.name,
		"--query", "Cluster.endpoint", "--output", "text", "--region", c.region).Output()
	if err != nil {
		log.Fatalf("Failed to get Cluster endpoint, %v", err)
	}
	c.clusterEndpoint = string(output)
}

func (c *Cluster) GetCert() {
	output, err := exec.Command("aws", "eks", "describe-cluster", "--name", c.name,
		"--query", "Cluster.certificateAuthority.data", "--output", "text", "--region", c.region).Output()
	if err != nil {
		log.Fatalf("Failed to get Cluster certificate data, %v", err)
	}
	c.certificateData = string(output)
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
