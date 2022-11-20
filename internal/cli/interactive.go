package cli

import (
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/manifoldco/promptui"
	"log"
)

func PromptUser() {
	searcher := func(input string, index int) bool {
		return fuzzy.Match(input, cfg.clusters[index].ClusterNickname) || fuzzy.Match(input, cfg.clusters[index].ClusterName)
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
		Items:     cfg.clusters,
		Size:      8,
		Templates: promptTemplate,
		Searcher:  searcher,
		IsVimMode: cfg.vimMode,
	}

	i, _, err := prompt.Run()
	if err != nil {
		log.Fatal(err)
	}

	cfg.clusters[i].SsoLogin()
	cfg.clusters[i].GetCert()
	cfg.clusters[i].GetEndpoint()
	cfg.clusters[i].GenerateKubeconfig()
	cfg.clusters[i].PrintExports()
}
