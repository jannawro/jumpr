package cli

import (
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/manifoldco/promptui"
	"log"
)

func PromptUser() {
	searcher := func(input string, index int) bool {
		return fuzzy.Match(input, cfg[index].Nickname()) || fuzzy.Match(input, cfg[index].Name())
	}

	promptTemplate := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "{{ .nickname | blue }} ({{ .clusterName | faint }})",
		Inactive: "{{ .nickname | faint }}",
		Selected: "Selected: {{ .nickname | blue }}",
		Details: `
--------- Cluster ----------
{{ "Nickname:" | faint }}	 {{ .nickname }}
{{ "Name:" | faint }}	 {{ .name }}
{{ "Region:" | faint }}	 {{ .region }}
{{ "Profile:" | faint }}	 {{ .profile }}`,
	}

	prompt := promptui.Select{
		Label:     "Select a cluster",
		Items:     cfg,
		Size:      8,
		Templates: promptTemplate,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()
	if err != nil {
		log.Fatal(err)
	}

	cfg[i].SsoLogin()
	cfg[i].GetCert()
	cfg[i].GetEndpoint()
	cfg[i].GenerateKubeconfig()
	cfg[i].PrintExports()
}
