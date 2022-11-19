package cli

import (
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/manifoldco/promptui"
	"log"
)

func PromptUser() {
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

	data[i].SsoLogin()
	data[i].GetCert()
	data[i].GetEndpoint()
	data[i].GenerateKubeconfig()
	data[i].PrintExports()
}
