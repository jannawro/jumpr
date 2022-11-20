apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: {{ .CertificateData }}
    server: {{ .ClusterEndpoint }}
  name: arn:aws:eks:{{ .AWSRegion }}:{{ .AWSAccountId }}:cluster/{{ .ClusterName }}
contexts:
- context:
    cluster: arn:aws:eks:{{ .AWSRegion }}:{{ .AWSAccountId }}:cluster/{{ .ClusterName }}
    user: arn:aws:eks:{{ .AWSRegion }}:{{ .AWSAccountId }}:cluster/{{ .ClusterName }}
  name: arn:aws:eks:{{ .AWSRegion }}:{{ .AWSAccountId }}:cluster/{{ .ClusterName }}
current-context: arn:aws:eks:{{ .AWSRegion }}:{{ .AWSAccountId }}:cluster/{{ .ClusterName }}
kind: Config
preferences: {}
users:
- name: arn:aws:eks:{{ .AWSRegion }}:{{ .AWSAccountId }}:cluster/{{ .ClusterName }}
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1beta1
      command: aws
      args:
        - --region
        - {{ .AWSRegion }}
        - eks
        - get-token
        - --cluster-name
        - {{ .ClusterName }}
