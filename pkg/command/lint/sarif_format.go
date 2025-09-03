package lint

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/template"

	"github.com/owenrumney/go-sarif/v2/sarif"
	"golang.stackrox.io/kube-linter/internal/consts"
	"golang.stackrox.io/kube-linter/pkg/command/checks"
	"golang.stackrox.io/kube-linter/pkg/command/common"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/run"
)

const (
	ruleHelpTemplateStr = `Check: {{.Name}}
Description: {{.Description}}
Remediation: {{.Remediation}}
Template: {{checkTemplateURL .}}`

	resultMessageTemplateStr = `{{.Report.Diagnostic.Message}}
object: {{.ObjectName}}`
)

var (
	ruleHelpTemplate = common.MustInstantiatePlainTemplate(ruleHelpTemplateStr,
		template.FuncMap{"checkTemplateURL": getCheckTemplateURL})

	resultMessageTemplate = common.MustInstantiatePlainTemplate(resultMessageTemplateStr, nil)
)

// formatLintSarif implements common.SARIFFormat.
// Must be used only with lint.Command because it only understands run.Result as data parameter.
func formatLintSarif(out io.Writer, data interface{}) error {
	if res, ok := data.(run.Result); ok {
		return formatSarif(out, res)
	}
	return errors.New("provided data must be of run.Result type")
}

func formatSarif(out io.Writer, result run.Result) error {
	sarifReport, err := sarif.New(sarif.Version210)
	if err != nil {
		return err
	}

	sarifRun := sarif.NewRunWithInformationURI(consts.ProgramName, consts.MainURL)
	sarifReport.AddRun(sarifRun)

	sarifRun.Tool.Driver.WithVersion(result.Summary.KubeLinterVersion)

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	sarifRun.AddInvocation(result.Summary.ChecksStatus == run.ChecksPassed).
		WithEndTimeUTC(result.Summary.CheckEndTime).
		// WithWorkingDirectory helps GitHub resolve artifact locations from repo root when their paths are absolute.
		WithWorkingDirectory(sarif.NewArtifactLocation().WithUri("file://" + cwd))

	for i := range result.Checks {
		err = addSarifRule(sarifRun, &result.Checks[i])
		if err != nil {
			return err
		}
	}

	for i := range result.Reports {
		err = addSarifResult(sarifRun, cwd, &result.Reports[i])
		if err != nil {
			return err
		}
	}

	return sarifReport.Write(out)
}

func addSarifRule(sarifRun *sarif.Run, check *config.Check) error {
	helpURL, err := getCheckTemplateURL(check)
	if err != nil {
		return err
	}

	helpText, err := renderTemplate(ruleHelpTemplate, check)
	if err != nil {
		return err
	}

	sarifRun.AddRule(check.Name).
		WithDescription(check.Description).
		WithFullDescription(sarif.NewMultiformatMessageString(check.Remediation)).
		WithHelpURI(helpURL).
		// Notice that we give WithHelp the same information as added above for couple of reasons:
		// 1) GitHub does not display HelpURI, although this attribute is required.
		// 2) Rule ID, short and full descriptions are shown at different spots on the screen but it is helpful to see
		//    them together.
		// Markdown format for Help seemed to be ignored therefore we only provide the plain text version.
		WithTextHelp(helpText)

	return nil
}

func getCheckTemplateURL(check *config.Check) (string, error) {
	anchor, err := checks.GetTemplateLink(check)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(consts.TemplateURLFormat, anchor), nil
}

func renderTemplate(t *template.Template, data interface{}) (string, error) {
	buf := &bytes.Buffer{}
	err := t.Execute(buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func addSarifResult(sarifRun *sarif.Run, cwd string, report *diagnostic.WithContext) error {
	sarifLocation := sarif.NewLocation()

	sarifLocation.PhysicalLocation = sarif.NewPhysicalLocation().
		WithArtifactLocation(sarif.NewArtifactLocation().WithUri(getArtifactURI(cwd, report.Object.Metadata.FilePath))).
		// We currently assign all errors to the first line in the file otherwise the absent region on the output
		// does not pass GitHub validation rule GH1003.
		// TODO: update region location with actual position in the file when it is available.
		WithRegion(sarif.NewRegion().WithStartLine(1))

	k8sObjectName := report.Object.GetK8sObjectName()

	// GitHub does not seem to show logical locations at the moment. We still provide them hoping it will in the future.
	sarifLocation.LogicalLocations = append(sarifLocation.LogicalLocations,
		sarif.NewLogicalLocation().WithName(k8sObjectName.Name).WithKind("Object Name"),
		sarif.NewLogicalLocation().WithName(k8sObjectName.Namespace).WithKind("Object Namespace"),
		sarif.NewLogicalLocation().WithName(k8sObjectName.GroupVersionKind.Group).WithKind("GVK/Group"),
		sarif.NewLogicalLocation().WithName(k8sObjectName.GroupVersionKind.Version).WithKind("GVK/Version").WithFullyQualifiedName(k8sObjectName.GroupVersionKind.GroupVersion().String()),
		sarif.NewLogicalLocation().WithName(k8sObjectName.GroupVersionKind.Kind).WithKind("GVK/Kind").WithFullyQualifiedName(k8sObjectName.GroupVersionKind.String()),
	)

	messageText, err := renderTemplate(resultMessageTemplate, struct {
		Report     *diagnostic.WithContext
		ObjectName lintcontext.K8sObjectInfo
	}{Report: report, ObjectName: k8sObjectName})
	if err != nil {
		return err
	}

	result := sarif.NewRuleResult(report.Check).
		WithMessage(sarif.NewTextMessage(messageText))
	result.AddLocation(sarifLocation)

	sarifRun.AddResult(result)

	return nil
}

// getArtifactURI tries to resolve path relative to cwd; if that fails, tries to get the absolute path with appended
// `file://` protocol; if that fails, returns the path as-is.
// GitHub prefers file URIs to be provided relative to the repo root. Assuming that this tool is invoked from the repo
// root, this function should resolve paths in a way GitHub likes them.
func getArtifactURI(cwd, path string) string {
	absolute, err := filepath.Abs(path)
	if err != nil {
		return path
	}

	relative, err := filepath.Rel(cwd, absolute)
	if err == nil {
		return relative
	}

	return "file://" + absolute
}
