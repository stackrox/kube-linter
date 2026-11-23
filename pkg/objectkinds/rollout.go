package objectkinds

import "k8s.io/apimachinery/pkg/runtime/schema"

var (
	rolloutGVK = schema.GroupVersion{Group: "argoproj.io", Version: "v1alpha1"}.WithKind("Rollout")
)
