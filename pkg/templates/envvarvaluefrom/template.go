package envvarvaluefrom

import (
	"fmt"
	"maps"
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

type resourceInfo struct {
	name     string
	key      string
	optional *bool
}

type resourceType int

const (
	resourceTypeSecret resourceType = iota
	resourceTypeConfigMap
)

type resourceChecker struct {
	objType      string
	objMap       map[string]interface{}
	getKeys      func(interface{}) []string
	ignoredRegex []*regexp.Regexp
}

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
					if secret, found := obj.K8sObject.(*v1.Secret); found {
						secrets[secret.Name] = secret // Fix: Remove ObjectMeta
					}
					if configmap, found := obj.K8sObject.(*v1.ConfigMap); found {
						configmaps[configmap.Name] = configmap // Fix: Remove ObjectMeta
					}
				}
				return lintForEachContainer(lintCtx, object, ignoredSecrets, ignoredConfigMaps, secrets, configmaps)
			}, nil
		}),
	})
}

func lintForEachContainer(lintCtx lintcontext.LintContext, object lintcontext.Object, ignoredSecrets, ignoredConfigMaps []*regexp.Regexp, secrets map[string]*v1.Secret, configmaps map[string]*v1.ConfigMap) []diagnostic.Diagnostic {
	return util.PerContainerCheck(func(container *v1.Container) []diagnostic.Diagnostic {
		var results []diagnostic.Diagnostic
		var envRefs []struct {
			info resourceInfo
			typ  resourceType
		}

		for _, envVar := range container.Env {
			valueFrom := envVar.ValueFrom
			if valueFrom == nil {
				continue
			}

			if secretRef := valueFrom.SecretKeyRef; secretRef != nil {
				envRefs = append(envRefs, struct {
					info resourceInfo
					typ  resourceType
				}{
					info: resourceInfo{
						name:     secretRef.Name,
						key:      secretRef.Key,
						optional: secretRef.Optional,
					},
					typ: resourceTypeSecret,
				})
			}

			if configMapRef := valueFrom.ConfigMapKeyRef; configMapRef != nil {
				envRefs = append(envRefs, struct {
					info resourceInfo
					typ  resourceType
				}{
					info: resourceInfo{
						name:     configMapRef.Name,
						key:      configMapRef.Key,
						optional: configMapRef.Optional,
					},
					typ: resourceTypeConfigMap,
				})
			}
		}

		secretChecker := &resourceChecker{
			objType:      "secret",
			objMap:       make(map[string]interface{}),
			getKeys:      getSecretKeys,
			ignoredRegex: ignoredSecrets,
		}
		for k, v := range secrets {
			secretChecker.objMap[k] = v
		}

		configMapChecker := &resourceChecker{
			objType:      "config map",
			objMap:       make(map[string]interface{}),
			getKeys:      getConfigMapKeys,
			ignoredRegex: ignoredConfigMaps,
		}
		for k, v := range configmaps {
			configMapChecker.objMap[k] = v
		}

		for _, envRef := range envRefs {
			var checker *resourceChecker
			switch envRef.typ {
			case resourceTypeSecret:
				checker = secretChecker
			case resourceTypeConfigMap:
				checker = configMapChecker
			}

			if msg := checkResourceReference(container.Name, envRef.info, checker); msg != "" {
				results = append(results, diagnostic.Diagnostic{Message: msg})
			}
		}
		return results
	})(lintCtx, object)
}

func checkResourceReference(containerName string, ref resourceInfo, checker *resourceChecker) string {
	if ref.optional != nil && *ref.optional {
		return ""
	}

	if isInRegexList(checker.ignoredRegex, ref.name) {
		return ""
	}

	obj, ok := checker.objMap[ref.name]
	if !ok {
		return fmt.Sprintf("The container %q is referring to an unknown %s %q", containerName, checker.objType, ref.name)
	}

	keys := checker.getKeys(obj)
	if !isInList(keys, ref.key) {
		return fmt.Sprintf("The container %q is referring to an unknown key %q in %s %q", containerName, ref.key, checker.objType, ref.name)
	}

	return ""
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

func getSecretKeys(obj interface{}) []string {
	secret := obj.(*v1.Secret)
	var keys []string
	for key := range maps.Keys(secret.Data) {
		keys = append(keys, key)
	}
	for key := range maps.Keys(secret.StringData) {
		keys = append(keys, key)
	}
	return keys
}

func getConfigMapKeys(obj interface{}) []string {
	configmap := obj.(*v1.ConfigMap)
	var keys []string
	for key := range maps.Keys(configmap.Data) {
		keys = append(keys, key)
	}
	for key := range maps.Keys(configmap.BinaryData) {
		keys = append(keys, key)
	}
	return keys
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
