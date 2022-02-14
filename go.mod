module golang.stackrox.io/kube-linter

go 1.16

require (
	github.com/Masterminds/sprig/v3 v3.2.2
	github.com/fatih/color v1.13.0
	github.com/ghodss/yaml v1.0.0
	github.com/golangci/golangci-lint v1.44.0
	github.com/mitchellh/mapstructure v1.4.3
	github.com/openshift/api v3.9.0+incompatible
	github.com/owenrumney/go-sarif/v2 v2.0.17
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.3.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.10.1
	github.com/stretchr/testify v1.7.0
	helm.sh/helm/v3 v3.7.0
	honnef.co/go/tools v0.2.2
	k8s.io/api v0.22.2
	k8s.io/apimachinery v0.22.2
	k8s.io/cli-runtime v0.22.2
	k8s.io/client-go v0.22.2
	k8s.io/gengo v0.0.0-20210915205010-39e73c8a59cd
)
