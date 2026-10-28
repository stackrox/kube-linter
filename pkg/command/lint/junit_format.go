package lint

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"

	"golang.stackrox.io/kube-linter/internal/consts"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/run"
)

type junitTestSuites struct {
	XMLName  xml.Name         `xml:"testsuites"`
	Name     string           `xml:"name,attr"`
	Tests    int              `xml:"tests,attr"`
	Failures int              `xml:"failures,attr"`
	Suites   []junitTestSuite `xml:"testsuite"`
}

type junitTestSuite struct {
	Name      string          `xml:"name,attr"`
	Tests     int             `xml:"tests,attr"`
	Failures  int             `xml:"failures,attr"`
	Timestamp string          `xml:"timestamp,attr,omitempty"`
	TestCases []junitTestCase `xml:"testcase"`
}

type junitTestCase struct {
	Name      string        `xml:"name,attr"`
	ClassName string        `xml:"classname,attr"`
	File      string        `xml:"file,attr,omitempty"`
	Failure   *junitFailure `xml:"failure,omitempty"`
}

type junitFailure struct {
	Message string `xml:"message,attr"`
	Type    string `xml:"type,attr"`
	Text    string `xml:",cdata"`
}

// formatLintJUnit implements common.JUnitFormat.
// Must be used only with lint.Command because it only understands run.Result as data parameter.
func formatLintJUnit(out io.Writer, data interface{}) error {
	if res, ok := data.(run.Result); ok {
		return formatJUnit(out, res)
	}
	return errors.New("provided data must be of run.Result type")
}

// formatJUnit writes one test suite per check. A check without lint errors is a single passing test case,
// and every lint error is a failing test case named after the object it was found on.
func formatJUnit(out io.Writer, result run.Result) error {
	reportsByCheck := make(map[string][]diagnostic.WithContext)
	var checkOrder []string
	for _, check := range result.Checks {
		if _, ok := reportsByCheck[check.Name]; !ok {
			reportsByCheck[check.Name] = nil
			checkOrder = append(checkOrder, check.Name)
		}
	}
	for _, report := range result.Reports {
		// Reports such as failed-to-load-object do not belong to an enabled check, so add their suite on first use.
		if _, ok := reportsByCheck[report.Check]; !ok {
			checkOrder = append(checkOrder, report.Check)
		}
		reportsByCheck[report.Check] = append(reportsByCheck[report.Check], report)
	}

	var timestamp string
	if !result.Summary.CheckEndTime.IsZero() {
		timestamp = result.Summary.CheckEndTime.UTC().Format("2006-01-02T15:04:05")
	}

	suites := junitTestSuites{Name: consts.ProgramName}
	for _, checkName := range checkOrder {
		suite := junitTestSuite{Name: checkName, Timestamp: timestamp}
		reports := reportsByCheck[checkName]
		if len(reports) == 0 {
			suite.TestCases = append(suite.TestCases, junitTestCase{Name: checkName, ClassName: checkName})
		}
		// Some JUnit consumers (GitLab, for example) keep only the first case with a given classname and name,
		// so repeated names in a suite get a counter.
		seen := make(map[string]int)
		for _, report := range reports {
			name := junitCaseName(report)
			seen[name]++
			if n := seen[name]; n > 1 {
				name = fmt.Sprintf("%s #%d", name, n)
			}
			suite.TestCases = append(suite.TestCases, junitTestCase{
				Name:      name,
				ClassName: checkName,
				File:      report.Object.Metadata.FilePath,
				Failure: &junitFailure{
					Message: report.Diagnostic.Message,
					Type:    checkName,
					Text:    junitFailureText(report),
				},
			})
		}
		suite.Tests = len(suite.TestCases)
		suite.Failures = len(reports)
		suites.Tests += suite.Tests
		suites.Failures += suite.Failures
		suites.Suites = append(suites.Suites, suite)
	}

	if _, err := io.WriteString(out, xml.Header); err != nil {
		return err
	}
	enc := xml.NewEncoder(out)
	enc.Indent("", "  ")
	if err := enc.Encode(suites); err != nil {
		return err
	}
	_, err := io.WriteString(out, "\n")
	return err
}

// junitCaseName names a failing case after its object. Objects that failed to load have no name, so their file
// path is used instead.
func junitCaseName(report diagnostic.WithContext) string {
	objectName := report.Object.GetK8sObjectName()
	if objectName.Name == "" && report.Object.Metadata.FilePath != "" {
		return report.Object.Metadata.FilePath
	}
	return objectName.String()
}

func junitFailureText(report diagnostic.WithContext) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s\n", report.Diagnostic.Message)
	fmt.Fprintf(&b, "object: %s\n", report.Object.GetK8sObjectName())
	if report.Object.Metadata.FilePath != "" {
		fmt.Fprintf(&b, "file: %s\n", report.Object.Metadata.FilePath)
	}
	if report.Remediation != "" {
		fmt.Fprintf(&b, "remediation: %s\n", report.Remediation)
	}
	return b.String()
}
