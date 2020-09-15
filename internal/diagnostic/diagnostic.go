package diagnostic

import (
	"fmt"

	"golang.stackrox.io/kube-linter/internal/lintcontext"
)

// A Diagnostic represents one specific problem diagnosed by a check.
type Diagnostic struct {
	Message string

	// TODO: add line number/col number
}

// WithContext puts a diagnostic in the context of which check emitted it,
// and which object it applied to.
type WithContext struct {
	Diagnostic Diagnostic
	Check      string
	Object     lintcontext.ObjectWithMetadata
}

// Format formats the diagnostic for human readability.
func (w *WithContext) Format() string {
	obj := w.Object.K8sObject
	return fmt.Sprintf("%s: %s (check: %s, object: %s %s/%s)", w.Object.FilePath, w.Diagnostic.Message, w.Check,
		obj.GetObjectKind().GroupVersionKind(), obj.GetNamespace(), obj.GetName())
}
