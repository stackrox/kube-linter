package envvarvaluefrom

import (
	"fmt"
	"maps"
	"regexp"

	"github.com/pkg/errors"
	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/extract"
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
	objects      map[string][]namespacedObject
	getKeys      func(interface{}) []string
	ignoredRegex []*regexp.Regexp
}

// namespacedObject pairs a secret or config map with its namespace: names repeat across namespaces.
type namespacedObject struct {
	namespace string
	object    interface{}
}

// inNamespace returns the objects a reference from namespace resolves against.
// An empty namespace matches any: manifests get one when they are applied.
func (c *resourceChecker) inNamespace(namespace, name string) []interface{} {
	var matches []interface{}
	for _, candidate := range c.objects[name] {
		if candidate.namespace == namespace || candidate.namespace == "" || namespace == "" {
			matches = append(matches, candidate.object)
		}
	}
	return matches
}

func init() {
	templates.Register(check.Template{
		HumanName:   "Env references",
		Key:         templateKey,
		Description: "Flag resources which use env variables or volumes from secrets/configmaps not included in the release",
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
				secrets := make(map[string][]namespacedObject)
				configmaps := make(map[string][]namespacedObject)
				for _, obj := range lintCtx.Objects() {
					if secret, found := obj.K8sObject.(*v1.Secret); found {
						secrets[secret.Name] = append(secrets[secret.Name], namespacedObject{namespace: secret.Namespace, object: secret})
					}
					if configmap, found := obj.K8sObject.(*v1.ConfigMap); found {
						configmaps[configmap.Name] = append(configmaps[configmap.Name], namespacedObject{namespace: configmap.Namespace, object: configmap})
					}
				}
				return lintForEachContainer(lintCtx, object, ignoredSecrets, ignoredConfigMaps, secrets, configmaps)
			}, nil
		}),
	})
}

func lintForEachContainer(lintCtx lintcontext.LintContext, object lintcontext.Object, ignoredSecrets, ignoredConfigMaps []*regexp.Regexp, secrets, configmaps map[string][]namespacedObject) []diagnostic.Diagnostic {
	namespace := object.K8sObject.GetNamespace()

	secretChecker := &resourceChecker{
		objType:      "secret",
		objects:      secrets,
		getKeys:      getSecretKeys,
		ignoredRegex: ignoredSecrets,
	}

	configMapChecker := &resourceChecker{
		objType:      "config map",
		objects:      configmaps,
		getKeys:      getConfigMapKeys,
		ignoredRegex: ignoredConfigMaps,
	}

	diagnostics := util.PerContainerCheck(func(container *v1.Container) []diagnostic.Diagnostic {
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

		for _, envFrom := range container.EnvFrom {
			if secretRef := envFrom.SecretRef; secretRef != nil {
				envRefs = append(envRefs, struct {
					info resourceInfo
					typ  resourceType
				}{
					info: resourceInfo{
						name:     secretRef.Name,
						optional: secretRef.Optional,
					},
					typ: resourceTypeSecret,
				})
			}

			if configMapRef := envFrom.ConfigMapRef; configMapRef != nil {
				envRefs = append(envRefs, struct {
					info resourceInfo
					typ  resourceType
				}{
					info: resourceInfo{
						name:     configMapRef.Name,
						optional: configMapRef.Optional,
					},
					typ: resourceTypeConfigMap,
				})
			}
		}

		subject := fmt.Sprintf("container %q", container.Name)
		for _, envRef := range envRefs {
			var checker *resourceChecker
			switch envRef.typ {
			case resourceTypeSecret:
				checker = secretChecker
			case resourceTypeConfigMap:
				checker = configMapChecker
			}

			if msg := checkResourceReference(subject, namespace, envRef.info, checker); msg != "" {
				results = append(results, diagnostic.Diagnostic{Message: msg})
			}
		}
		return results
	})(lintCtx, object)

	return append(diagnostics, lintVolumes(object, namespace, secretChecker, configMapChecker)...)
}

// lintVolumes checks the secrets and config maps mounted as pod volumes.
func lintVolumes(object lintcontext.Object, namespace string, secretChecker, configMapChecker *resourceChecker) []diagnostic.Diagnostic {
	podSpec, found := extract.PodSpec(object.K8sObject)
	if !found {
		return nil
	}

	var results []diagnostic.Diagnostic
	check := func(subject string, ref resourceInfo, items []v1.KeyToPath, checker *resourceChecker) {
		// Check the object itself first so that a missing object is reported once, not once per item.
		refs := []resourceInfo{ref}
		for _, item := range items {
			refs = append(refs, resourceInfo{name: ref.name, key: item.Key, optional: ref.optional})
		}
		for _, r := range refs {
			if msg := checkResourceReference(subject, namespace, r, checker); msg != "" {
				results = append(results, diagnostic.Diagnostic{Message: msg})
				if r.key == "" {
					return
				}
			}
		}
	}

	for _, volume := range podSpec.Volumes {
		subject := fmt.Sprintf("volume %q", volume.Name)
		if secret := volume.Secret; secret != nil {
			check(subject, resourceInfo{name: secret.SecretName, optional: secret.Optional}, secret.Items, secretChecker)
		}
		if configMap := volume.ConfigMap; configMap != nil {
			check(subject, resourceInfo{name: configMap.Name, optional: configMap.Optional}, configMap.Items, configMapChecker)
		}
	}
	return results
}

// checkResourceReference returns a message if ref points to a missing object or key.
// An empty key only requires the object to exist.
func checkResourceReference(subject, namespace string, ref resourceInfo, checker *resourceChecker) string {
	if ref.optional != nil && *ref.optional {
		return ""
	}

	if isInRegexList(checker.ignoredRegex, ref.name) {
		return ""
	}

	objs := checker.inNamespace(namespace, ref.name)
	if len(objs) == 0 {
		return fmt.Sprintf("The %s is referring to an unknown %s %q", subject, checker.objType, ref.name)
	}

	if ref.key == "" {
		return ""
	}

	for _, obj := range objs {
		if isInList(checker.getKeys(obj), ref.key) {
			return ""
		}
	}

	return fmt.Sprintf("The %s is referring to an unknown key %q in %s %q", subject, ref.key, checker.objType, ref.name)
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
