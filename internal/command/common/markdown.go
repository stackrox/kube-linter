package common

import (
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"golang.stackrox.io/kube-linter/internal/utils"
)

// MustInstantiateTemplate instanties the given go template with a common list of
// functions. It panics if there is an error.
func MustInstantiateTemplate(templateStr string) *template.Template {
	tpl, err := template.New("").Funcs(sprig.TxtFuncMap()).Funcs(
		template.FuncMap{
			"backtick": func() string {
				return "`"
			},
		},
	).Parse(templateStr)
	utils.Must(err)
	return tpl

}
