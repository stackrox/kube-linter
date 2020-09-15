package diagnostic

import (
	"fmt"
	"io"

	"github.com/fatih/color"

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

var (
	bold = color.New(color.Bold)
)

// FormatTo formats the diagnostic for human readability.
func (w *WithContext) FormatTo(out io.Writer) {
	fmt.Fprintf(out, "%s %s", bold.Sprintf("%s:", w.Object.FilePath), color.RedString(w.Diagnostic.Message))
	obj := w.Object.K8sObject
	fmt.Fprintf(out, " (check: %s, object: %s %s)\n\n", color.YellowString(w.Check),
		color.YellowString("%s/%s", obj.GetNamespace(), obj.GetName()), color.YellowString(obj.GetObjectKind().GroupVersionKind().String()))

}
