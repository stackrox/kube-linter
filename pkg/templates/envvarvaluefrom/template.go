package envvarvaluefrom

import (
	"fmt"
	"regexp"

	"github.com/pkg/errors"
	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/envvarvaluefrom/internal/params"
	"golang.stackrox.io/kube-linter/pkg/templates/util"
	v1 "k8s.io/api/core/v1"
)

const (
	templateKey = "env-value-from"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "Env references",
		Key:         templateKey,
		Description: "Flag resources which use env variables from secrets/configmaps not included in the release",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.DeploymentLike},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(p params.Params) (check.Func, error) {
			ignoredSecrets, err := extractRegexList(p.IgnoredSecrets)
			if err != nil {
				return nil, err
			}

			ignoredConfigMaps, err := extractRegexList(p.IgnoredConfigMaps)
			if err != nil {
				return nil, err
			}

			return func(lintCtx lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
				secrets := make(map[string]*v1.Secret)
				configmaps := make(map[string]*v1.ConfigMap)

				for _, obj := range lintCtx.Objects() {
					secret, found := obj.K8sObject.(*v1.Secret)
					if found {
						secrets[secret.ObjectMeta.Name] = secret
					}

					configmap, found := obj.K8sObject.(*v1.ConfigMap)
					if found {
						configmaps[configmap.ObjectMeta.Name] = configmap
					}
				}

				return lintForEachContainer(lintCtx, object, ignoredSecrets, ignoredConfigMaps, secrets, configmaps)
			}, nil
		}),
	})
}

func lintForEachContainer(lintCtx lintcontext.LintContext, object lintcontext.Object, ignoredSecrets []*regexp.Regexp, ignoredConfigMaps []*regexp.Regexp, secrets map[string]*v1.Secret, configmaps map[string]*v1.ConfigMap) []diagnostic.Diagnostic {
	return util.PerContainerCheck(func(container *v1.Container) []diagnostic.Diagnostic {
		var results []diagnostic.Diagnostic
		for _, envVar := range container.Env {
			valueFrom := envVar.ValueFrom
			if valueFrom == nil {
				continue
			}

			if secretKeySelector := valueFrom.SecretKeyRef; secretKeySelector != nil {
				if secretKeySelector.Optional != nil && *secretKeySelector.Optional {
					continue
				}

				if len(ignoredSecrets) > 0 && isInRegexList(ignoredSecrets, secretKeySelector.Name) {
					continue
				}

				secret, ok := secrets[secretKeySelector.Name]
				if !ok {
					results = append(results, diagnostic.Diagnostic{
						Message: fmt.Sprintf("The container %q is referring to an unknown secret %q", container.Name, secretKeySelector.Name),
					})
					continue
				}
				if isInList(Keys(secret.Data), secretKeySelector.Key) || isInList(Keys(secret.StringData), secretKeySelector.Key) {
					continue
				}
				results = append(results, diagnostic.Diagnostic{
					Message: fmt.Sprintf("The container %q is referring to an unknown key %q in secret %q", container.Name, secretKeySelector.Key, secretKeySelector.Name),
				})
			}

			if configMapSelector := valueFrom.ConfigMapKeyRef; configMapSelector != nil {
				if configMapSelector.Optional != nil && *configMapSelector.Optional {
					continue
				}
				if len(ignoredConfigMaps) > 0 && isInRegexList(ignoredConfigMaps, configMapSelector.Name) {
					continue
				}

				configmap, ok := configmaps[configMapSelector.Name]
				if !ok {
					results = append(results, diagnostic.Diagnostic{
						Message: fmt.Sprintf("The container %q is referring to an unknown config map %q", container.Name, configMapSelector.Name),
					})
					continue
				}

				if isInList(Keys(configmap.Data), configMapSelector.Key) || isInList(Keys(configmap.BinaryData), configMapSelector.Key) {
					continue
				}

				results = append(results, diagnostic.Diagnostic{
					Message: fmt.Sprintf("The container %q is referring to an unknown key %q in config map %q", container.Name, configMapSelector.Key, configMapSelector.Name),
				})
			}
		}
		return results
	})(lintCtx, object)
}

func isInRegexList(regexlist []*regexp.Regexp, name string) bool {
	for _, regex := range regexlist {
		if regex.MatchString(name) {
			return true
		}
	}
	return false
}

func isInList(regexlist []string, name string) bool {
	for _, regex := range regexlist {
		if name == regex {
			return true
		}
	}
	return false
}

func Keys[K comparable, V any](m map[K]V) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	return r
}

func extractRegexList(inputList []string) ([]*regexp.Regexp, error) {
	result := make([]*regexp.Regexp, 0, len(inputList))
	for _, res := range inputList {
		rg, err := regexp.Compile(res)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid regex %s", res)
		}
		result = append(result, rg)
	}
	return result, nil
}
