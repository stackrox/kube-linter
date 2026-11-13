package convert

import (
	"fmt"

	"github.com/ugiordan/kube-chainsaw/pkg/models"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	appsV1 "k8s.io/api/apps/v1"
	batchV1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	rbacV1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// FromLintContext converts kube-linter's typed K8s objects into
// kube-chainsaw's data structures. The Doc fields carry the full
// unstructured objects because kube-chainsaw's analyzer reads pod specs,
// aggregation rules and other content straight from Doc.
func FromLintContext(ctx lintcontext.LintContext) (*models.LoadedResources, error) {
	resources := models.NewLoadedResources()

	for _, obj := range ctx.Objects() {
		path := obj.Metadata.FilePath
		k8sObj := obj.K8sObject
		if k8sObj == nil {
			return nil, fmt.Errorf("nil K8sObject at %s", path)
		}

		doc, err := runtime.DefaultUnstructuredConverter.ToUnstructured(k8sObj)
		if err != nil {
			return nil, fmt.Errorf("converting object at %s to unstructured: %w", path, err)
		}

		switch o := k8sObj.(type) {
		case *rbacV1.ClusterRole:
			if existing, exists := resources.ClusterRoles[o.Name]; exists {
				return nil, duplicateError("ClusterRole", o.Name, path, existing.File)
			}
			resources.ClusterRoles[o.Name] = convertClusterRole(path, doc)
		case *rbacV1.Role:
			key := DefaultNamespace("Role", o.Namespace) + "/" + o.Name
			if existing, exists := resources.Roles[key]; exists {
				return nil, duplicateError("Role", key, path, existing.File)
			}
			resources.Roles[key] = convertRole(o, path, doc)
		case *rbacV1.ClusterRoleBinding:
			resources.ClusterRoleBindings = append(resources.ClusterRoleBindings, convertBinding(o.Name, "", o.RoleRef, o.Subjects, path, doc))
		case *rbacV1.RoleBinding:
			resources.RoleBindings = append(resources.RoleBindings, convertBinding(o.Name, DefaultNamespace("RoleBinding", o.Namespace), o.RoleRef, o.Subjects, path, doc))
		case *v1.ServiceAccount:
			namespace := DefaultNamespace("ServiceAccount", o.Namespace)
			key := namespace + "/" + o.Name
			if existing, exists := resources.ServiceAccounts[key]; exists {
				return nil, duplicateError("ServiceAccount", key, path, existing.File)
			}
			resources.ServiceAccounts[key] = &models.SAData{
				Name: o.Name, Namespace: namespace, File: path, Doc: doc,
			}
		case *v1.Pod:
			namespace := DefaultNamespace("Pod", o.Namespace)
			key := namespace + "/" + o.Name
			if existing, exists := resources.Pods[key]; exists {
				return nil, duplicateError("Pod", key, path, existing.File)
			}
			resources.Pods[key] = &models.PodData{
				Name: o.Name, Namespace: namespace,
				ServiceAccountName: serviceAccountName(o.Spec), File: path,
				Doc: doc,
			}
		case *appsV1.Deployment:
			if err := convertWorkload(resources, "Deployment", o.Name, o.Namespace, o.Spec.Template.Spec, path, doc); err != nil {
				return nil, err
			}
		case *appsV1.DaemonSet:
			if err := convertWorkload(resources, "DaemonSet", o.Name, o.Namespace, o.Spec.Template.Spec, path, doc); err != nil {
				return nil, err
			}
		case *appsV1.StatefulSet:
			if err := convertWorkload(resources, "StatefulSet", o.Name, o.Namespace, o.Spec.Template.Spec, path, doc); err != nil {
				return nil, err
			}
		case *appsV1.ReplicaSet:
			if err := convertWorkload(resources, "ReplicaSet", o.Name, o.Namespace, o.Spec.Template.Spec, path, doc); err != nil {
				return nil, err
			}
		case *batchV1.Job:
			if err := convertWorkload(resources, "Job", o.Name, o.Namespace, o.Spec.Template.Spec, path, doc); err != nil {
				return nil, err
			}
		case *batchV1.CronJob:
			if err := convertWorkload(resources, "CronJob", o.Name, o.Namespace, o.Spec.JobTemplate.Spec.Template.Spec, path, doc); err != nil {
				return nil, err
			}
		}
	}

	return resources, nil
}

// DefaultNamespace returns the namespace kube-chainsaw identifies a namespaced
// object by, defaulting an omitted namespace to "default" the same way its
// loader does. Cluster-scoped kinds keep the empty namespace.
func DefaultNamespace(kind, namespace string) string {
	if namespace == "" && isNamespacedKind(kind) {
		return "default"
	}
	return namespace
}

func isNamespacedKind(kind string) bool {
	switch kind {
	case "Role", "RoleBinding", "ServiceAccount", "Pod",
		"Deployment", "DaemonSet", "StatefulSet", "Job", "CronJob", "ReplicaSet":
		return true
	}
	return false
}

// serviceAccountName resolves the service account a pod spec runs as, falling
// back to the deprecated serviceAccount field before defaulting, matching the
// API server's behavior.
func serviceAccountName(podSpec v1.PodSpec) string {
	if podSpec.ServiceAccountName != "" {
		return podSpec.ServiceAccountName
	}
	if podSpec.DeprecatedServiceAccount != "" {
		return podSpec.DeprecatedServiceAccount
	}
	return "default"
}

func duplicateError(kind, name, path, existingPath string) error {
	return fmt.Errorf("duplicate %s %q at %s (already defined at %s)", kind, name, path, existingPath)
}

func convertClusterRole(path string, doc map[string]interface{}) *models.ClusterRoleData {
	return &models.ClusterRoleData{Rules: rulesFromDoc(doc), File: path, Doc: doc}
}

func convertRole(r *rbacV1.Role, path string, doc map[string]interface{}) *models.RoleData {
	return &models.RoleData{
		Rules: rulesFromDoc(doc), Namespace: DefaultNamespace("Role", r.Namespace), File: path, Doc: doc,
	}
}

// rulesFromDoc extracts RBAC rules from the object's unstructured document,
// mirroring kube-chainsaw's own loader. This keeps every rule field intact
// instead of reconstructing a subset that has to be kept in sync with the
// analyzer.
func rulesFromDoc(doc map[string]interface{}) []map[string]interface{} {
	rulesRaw, ok := doc["rules"]
	if !ok {
		return nil
	}
	rulesSlice, ok := rulesRaw.([]interface{})
	if !ok {
		return nil
	}
	var rules []map[string]interface{}
	for _, r := range rulesSlice {
		if rule, ok := r.(map[string]interface{}); ok {
			rules = append(rules, rule)
		}
	}
	return rules
}

func convertBinding(name, namespace string, roleRef rbacV1.RoleRef, subjects []rbacV1.Subject, path string, doc map[string]interface{}) *models.BindingData {
	subs := make([]map[string]interface{}, len(subjects))
	for i, s := range subjects {
		subs[i] = map[string]interface{}{
			"kind":      s.Kind,
			"name":      s.Name,
			"namespace": s.Namespace,
		}
	}
	return &models.BindingData{
		Name:      name,
		Namespace: namespace,
		RoleRef: map[string]interface{}{
			"kind": roleRef.Kind,
			"name": roleRef.Name,
		},
		Subjects: subs,
		File:     path,
		Doc:      doc,
	}
}

func convertWorkload(resources *models.LoadedResources, kind, name, namespace string, podSpec v1.PodSpec, path string, doc map[string]interface{}) error {
	namespace = DefaultNamespace(kind, namespace)
	key := kind + "/" + namespace + "/" + name
	if existing, exists := resources.Workloads[key]; exists {
		return duplicateError(kind, key, path, existing.File)
	}
	resources.Workloads[key] = &models.WorkloadData{
		Name: name, Kind: kind, Namespace: namespace,
		ServiceAccountName: serviceAccountName(podSpec), File: path,
		Doc: doc,
	}
	return nil
}
