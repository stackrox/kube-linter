package diagnostic

import (
	"bytes"
	"fmt"
	"io"

	"github.com/acarl005/stripansi"
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
	Remediation string
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
	fmt.Fprintf(out, "%s %s",
		bold.Sprintf("%s: (object: %s)", w.Object.Metadata.FilePath, formatObj(w.Object.K8sObject)), color.RedString(w.Diagnostic.Message),
	)
	fmt.Fprintf(out, " (check: %s, remediation: %s)\n\n", color.YellowString(w.Check),
		color.YellowString(w.Remediation))
}

// FormatPlain prints out the result to the given writer, without colors/special formatting.
func (w *WithContext) FormatPlain(out io.Writer) {
	var buf bytes.Buffer
	w.FormatToTerminal(&buf)
	out.Write([]byte(stripansi.Strip(buf.String())))
}
