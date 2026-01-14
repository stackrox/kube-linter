package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// AddDirectlyToScheme adds StackRox types to the scheme.
func AddDirectlyToScheme(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(GroupVersion,
		&Central{},
		&CentralList{},
		&SecuredCluster{},
		&SecuredClusterList{},
	)
	metav1.AddToGroupVersion(scheme, GroupVersion)
	return nil
}
