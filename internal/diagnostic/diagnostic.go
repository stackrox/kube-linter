package diagnostic

import (
	"fmt"
	"io"

	"github.com/fatih/color"
	"golang.stackrox.io/kube-linter/internal/k8sutil"
	"golang.stackrox.io/kube-linter/internal/lintcontext"
	"golang.stackrox.io/kube-linter/internal/stringutils"
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
	Object     lintcontext.Object
}

var (
	bold = color.New(color.Bold)
)

func formatObj(obj k8sutil.Object) string {
	return fmt.Sprintf("%s/%s %s", stringutils.OrDefault(obj.GetNamespace(), "<no namespace>"), obj.GetName(), obj.GetObjectKind().GroupVersionKind())
}

// FormatToTerminal writes the result to the given writer, which is expected to support
// terminal-based formatting.
func (w *WithContext) FormatToTerminal(out io.Writer) {
	fmt.Fprintf(out, "%s %s", bold.Sprintf("%s:", w.Object.Metadata.FilePath), color.RedString(w.Diagnostic.Message))
	fmt.Fprintf(out, " (check: %s, object: %s)\n\n", color.YellowString(w.Check),
		color.YellowString(formatObj(w.Object.K8sObject)))
}

// FormatPlain prints out the result to the given writer, without colors/special formatting.
func (w *WithContext) FormatPlain(out io.Writer) {
	fmt.Fprintf(out, "%s %s (check: %s, object: %s)\n\n", w.Object.Metadata.FilePath, w.Diagnostic.Message, w.Check, formatObj(w.Object.K8sObject))
}
