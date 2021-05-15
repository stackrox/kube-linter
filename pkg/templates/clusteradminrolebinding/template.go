package clusteradminrolebinding

import (
	"encoding/json"
	"fmt"

	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/clusteradminrolebinding/internal/params"
	rbacV1 "k8s.io/api/rbac/v1"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "cluster-admin Role Binding",
		Key:         "cluster-admin-role-binding",
		Description: "Flag bindings of cluster-admin role to service accounts, users, or groups",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.ClusterRoleBinding},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(_ params.Params) (check.Func, error) {
			return func(lintCtx lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
				clusterrolebinding, ok := object.K8sObject.(*rbacV1.ClusterRoleBinding)
				if !ok {
					return nil
				}
				clusterrole := clusterrolebinding.RoleRef
				if clusterrole.Name == "cluster-admin" && clusterrole.Kind == "ClusterRole" {
					jsonObj, _ := json.Marshal(clusterrolebinding.Subjects)
					return []diagnostic.Diagnostic{{Message: fmt.Sprintf("%s role is bound to %v", clusterrole.Name, string(jsonObj))}}
				}
				return nil
			}, nil
		}),
	})
}
