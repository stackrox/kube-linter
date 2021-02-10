package common

import (
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"golang.stackrox.io/kube-linter/internal/utils"
)

// MustInstantiateTemplate instantiates the given go template with a common list of
// functions. It panics if there is an error.
func MustInstantiateTemplate(templateStr string, customFuncMap template.FuncMap) *template.Template {
	tpl, err := template.New("").Funcs(sprig.TxtFuncMap()).Funcs(
		template.FuncMap{
			"codeSnippet": func(code string) string {
				return "`" + code + "`"
			},
			"codeSnippetInTable": func(code string) string {
				return "`" + strings.ReplaceAll(code, "|", `\|`) + "`"
			},
			"codeBlock": func(lang, code string) string {
				finalNewline := "\n"
				if strings.HasSuffix(code, "\n") {
					finalNewline = ""
				}
				return "```" + lang + "\n" + code + finalNewline + "```"
			},
		},
	).Funcs(customFuncMap).Parse(templateStr)
	utils.Must(err)
	return tpl

}
