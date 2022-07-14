package common_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.stackrox.io/kube-linter/pkg/command/common"
)

func TestMarkdownFunctions(t *testing.T) {
	templateTests := []struct {
		name string
		in   string
		data string
		out  string
	}{
		{"template without replacements", "foo bar baz", "", "foo bar baz"},
		{"template with replacement without filter", "foo bar {{ . }} baz", "code snippet", "foo bar code snippet baz"},
		{"template with replacement with codeSnippet filter", "foo bar {{ . | codeSnippet }} baz", "code|snippet", "foo bar `code|snippet` baz"},
		{"template with replacement with codeSnippetInTable filter (no escaping)", "foo bar {{ . | codeSnippetInTable }} baz", "code := snippet()", "foo bar `code := snippet()` baz"},
		{"template with replacement with codeSnippetInTable filter (escaping)", "foo bar {{ . | codeSnippetInTable }} baz", "code|snippet", "foo bar `code\\|snippet` baz"},
		{"template with replacement with codeBlock filter", "foo bar\n{{ . | codeBlock \"\" }}\nbaz", "code|snippet", "foo bar\n```\ncode|snippet\n```\nbaz"},
		{"template with replacement with codeBlock filter with lang", "foo bar\n{{ . | codeBlock \"txt\" }}\nbaz", "code|snippet", "foo bar\n```txt\ncode|snippet\n```\nbaz"},
	}

	for _, tt := range templateTests {
		t.Run(tt.name, func(t *testing.T) {
			tpl := common.MustInstantiateMarkdownTemplate(tt.in, nil)

			var b bytes.Buffer

			require.NoError(t, tpl.Execute(&b, tt.data))

			assert.Equal(t, tt.out, b.String())
		})
	}
}
