package convert

import (
	"fmt"

	"github.com/ugiordan/kube-chainsaw/pkg/models"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	appsV1 "k8s.io/api/apps/v1"
	batchV1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	rbacV1 "k8s.io/api/rbac/v1"
)

// FromLintContext converts kube-linter's typed K8s objects into
// kube-chainsaw's map[string]interface{} data structures.
func FromLintContext(ctx lintcontext.LintContext) (*models.LoadedResources, error) {
	resources := models.NewLoadedResources()

	for _, obj := range ctx.Objects() {
		path := obj.Metadata.FilePath
		k8sObj := obj.K8sObject

		switch o := k8sObj.(type) {
		case *rbacV1.ClusterRole:
			resources.ClusterRoles[o.Name] = convertClusterRole(o, path)
		case *rbacV1.Role:
			key := o.Namespace + "/" + o.Name
			resources.Roles[key] = convertRole(o, path)
		case *rbacV1.ClusterRoleBinding:
			resources.ClusterRoleBindings = append(resources.ClusterRoleBindings, convertBinding("ClusterRoleBinding", o.Name, "", o.RoleRef, o.Subjects, path))
		case *rbacV1.RoleBinding:
			resources.RoleBindings = append(resources.RoleBindings, convertBinding("RoleBinding", o.Name, o.Namespace, o.RoleRef, o.Subjects, path))
		case *v1.ServiceAccount:
			key := o.Namespace + "/" + o.Name
			resources.ServiceAccounts[key] = &models.SAData{
				Name: o.Name, Namespace: o.Namespace, File: path,
				Doc: map[string]interface{}{"kind": "ServiceAccount", "metadata": map[string]interface{}{"name": o.Name, "namespace": o.Namespace}},
			}
		case *v1.Pod:
			key := o.Namespace + "/" + o.Name
			saName := o.Spec.ServiceAccountName
			if saName == "" {
				saName = "default"
			}
			resources.Pods[key] = &models.PodData{
				Name: o.Name, Namespace: o.Namespace,
				ServiceAccountName: saName, File: path,
				Doc: map[string]interface{}{"kind": "Pod", "metadata": map[string]interface{}{"name": o.Name, "namespace": o.Namespace}},
			}
		case *appsV1.Deployment:
			convertWorkload(resources, "Deployment", o.Name, o.Namespace, o.Spec.Template.Spec, path)
		case *appsV1.DaemonSet:
			convertWorkload(resources, "DaemonSet", o.Name, o.Namespace, o.Spec.Template.Spec, path)
		case *appsV1.StatefulSet:
			convertWorkload(resources, "StatefulSet", o.Name, o.Namespace, o.Spec.Template.Spec, path)
		case *appsV1.ReplicaSet:
			convertWorkload(resources, "ReplicaSet", o.Name, o.Namespace, o.Spec.Template.Spec, path)
		case *batchV1.Job:
			convertWorkload(resources, "Job", o.Name, o.Namespace, o.Spec.Template.Spec, path)
		case *batchV1.CronJob:
			convertWorkload(resources, "CronJob", o.Name, o.Namespace, o.Spec.JobTemplate.Spec.Template.Spec, path)
		}
	}

	return resources, nil
}

func convertClusterRole(cr *rbacV1.ClusterRole, path string) *models.ClusterRoleData {
	rules := make([]map[string]interface{}, len(cr.Rules))
	for i, r := range cr.Rules {
		rules[i] = convertPolicyRule(r)
	}
	doc := map[string]interface{}{
		"kind":     "ClusterRole",
		"metadata": map[string]interface{}{"name": cr.Name},
	}
	if cr.AggregationRule != nil {
		doc["aggregationRule"] = map[string]interface{}{
			"clusterRoleSelectors": fmt.Sprintf("%v", cr.AggregationRule.ClusterRoleSelectors),
		}
	}
	return &models.ClusterRoleData{Rules: rules, File: path, Doc: doc}
}

func convertRole(r *rbacV1.Role, path string) *models.RoleData {
	rules := make([]map[string]interface{}, len(r.Rules))
	for i, rule := range r.Rules {
		rules[i] = convertPolicyRule(rule)
	}
	return &models.RoleData{
		Rules: rules, Namespace: r.Namespace, File: path,
		Doc: map[string]interface{}{
			"kind":     "Role",
			"metadata": map[string]interface{}{"name": r.Name, "namespace": r.Namespace},
		},
	}
}

func convertPolicyRule(r rbacV1.PolicyRule) map[string]interface{} {
	return map[string]interface{}{
		"apiGroups": toInterfaceSlice(r.APIGroups),
		"resources": toInterfaceSlice(r.Resources),
		"verbs":     toInterfaceSlice(r.Verbs),
	}
}

func convertBinding(kind, name, namespace string, roleRef rbacV1.RoleRef, subjects []rbacV1.Subject, path string) *models.BindingData {
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
		Doc: map[string]interface{}{
			"kind":     kind,
			"metadata": map[string]interface{}{"name": name, "namespace": namespace},
		},
	}
}

func convertWorkload(resources *models.LoadedResources, kind, name, namespace string, podSpec v1.PodSpec, path string) {
	saName := podSpec.ServiceAccountName
	if saName == "" {
		saName = "default"
	}
	key := kind + "/" + namespace + "/" + name
	resources.Workloads[key] = &models.WorkloadData{
		Name: name, Kind: kind, Namespace: namespace,
		ServiceAccountName: saName, File: path,
		Doc: map[string]interface{}{"kind": kind, "metadata": map[string]interface{}{"name": name, "namespace": namespace}},
	}
}

func toInterfaceSlice(ss []string) []interface{} {
	result := make([]interface{}, len(ss))
	for i, s := range ss {
		result[i] = s
	}
	return result
}
