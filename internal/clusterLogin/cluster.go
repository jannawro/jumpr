package clusterLogin

import (
	"bytes"
	"context"
	"embed"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eks"
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
	AwsConfig       aws.Config
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

	c.AwsConfig, err = config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(c.Profile))
	if err != nil {
		log.Fatalf("Failed loading aws profile from ~/.aws/config: %v", err)
	}
}

func (c *Cluster) GetClusterInfo() {
	client := eks.NewFromConfig(c.AwsConfig, func(o *eks.Options) {
		o.Region = c.Region
	})

	resp, err := client.DescribeCluster(context.TODO(), &eks.DescribeClusterInput{
		Name: aws.String(c.Name),
	})
	if err != nil {
		log.Fatal("Failed at getting cluster information:\n", err)
	}

	c.ClusterEndpoint = aws.ToString(resp.Cluster.Endpoint)
	c.CertificateData = aws.ToString(resp.Cluster.CertificateAuthority.Data)
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
