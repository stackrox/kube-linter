package common

import (
	"strings"
	"text/template"

	"github.com/fatih/color"

	"github.com/Masterminds/sprig/v3"
	"golang.stackrox.io/kube-linter/internal/utils"
)

var (
	colorRed    = color.New(color.FgRed)
	colorYellow = color.New(color.FgYellow)
	colorBold   = color.New(color.Bold)

	markdownFuncs = template.FuncMap{
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
	}

	plainFuncs = template.FuncMap{
		"red":    colorRed.Sprint,
		"yellow": colorYellow.Sprint,
		"bold":   colorBold.Sprint,
	}
)

// MustInstantiateMarkdownTemplate instantiates the given go template with a common list of markdown functions.
// It panics if there is an error.
func MustInstantiateMarkdownTemplate(templateStr string, customFuncMap template.FuncMap) *template.Template {
	tpl, err := instantiateTemplate(templateStr, markdownFuncs, customFuncMap)
	utils.Must(err)
	return tpl
}

// MustInstantiatePlainTemplate instantiates the given go template with a common list of functions for console output.
// It panics if there is an error.
func MustInstantiatePlainTemplate(templateStr string, customFuncMap template.FuncMap) *template.Template {
	tpl, err := instantiateTemplate(templateStr, plainFuncs, customFuncMap)
	utils.Must(err)
	return tpl
}

func instantiateTemplate(templateStr string, commonFuncMap, customFuncMap template.FuncMap) (*template.Template, error) {
	tpl, err := template.New("").Funcs(sprig.TxtFuncMap()).Funcs(commonFuncMap).Funcs(customFuncMap).Parse(templateStr)
	return tpl, err
}
