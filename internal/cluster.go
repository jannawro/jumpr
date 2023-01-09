package internal

import (
	"bytes"
	"context"
	"embed"
	"os"
	"os/exec"
	"text/template"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eks"
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
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	check("Failed completing the login process:", err)

	c.AwsConfig, err = config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(c.Profile))
	check("Failed loading aws profile from ~/.aws/config:", err)
}

func (c *Cluster) GetClusterInfo() {
	client := eks.NewFromConfig(c.AwsConfig, func(o *eks.Options) {
		o.Region = c.Region
	})

	resp, err := client.DescribeCluster(context.TODO(), &eks.DescribeClusterInput{
		Name: aws.String(c.Name),
	})
	check("Failed at getting cluster information:", err)

	c.ClusterEndpoint = aws.ToString(resp.Cluster.Endpoint)
	c.CertificateData = aws.ToString(resp.Cluster.CertificateAuthority.Data)
}

func (c *Cluster) GenerateKubeconfig() {
	tmpl, err := template.New("kubeconfig.gotmpl").ParseFS(templates, "templates/kubeconfig.gotmpl")
	check("Failed to parse kubeconfig template:", err)

	b := new(bytes.Buffer)

	err = tmpl.Execute(b, c)
	check("Failed to execute kubeconfig template:", err)

	kubePath := "/tmp/kubeconfig-" + c.Name

	c.KubeconfigPath = kubePath

	f, err := os.Create(kubePath)
	check("Failed creating a file in /tmp/:", err)

	defer func() {
		err = f.Close()
		check("Failed at closing the file kubeconfig file:", err)
	}()

	_, err = f.WriteString(b.String())
}

func (c *Cluster) PrintExports() {
	tmpl, err := template.New("exports.gotmpl").ParseFS(templates, "templates/exports.gotmpl")
	check("Failed to parse kubeconfig template:", err)

	err = tmpl.Execute(os.Stdout, c)
	check("Failed to execute exports template:", err)
}
