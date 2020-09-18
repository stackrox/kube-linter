package common

import (
	"strings"
	"text/template"

	"golang.stackrox.io/kube-linter/internal/utils"
)

// MustInstantiateTemplate instanties the given go template with a common list of
// functions. It panics if there is an error.
func MustInstantiateTemplate(templateStr string) *template.Template {
	tpl, err := template.New("").Funcs(
		template.FuncMap{
			"backtick": func() string {
				return "`"
			},
			"joinstrings": strings.Join,
		},
	).Parse(templateStr)
	utils.Must(err)
	return tpl

}
