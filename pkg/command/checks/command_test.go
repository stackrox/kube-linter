package checks

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.stackrox.io/kube-linter/internal/set"
	"golang.stackrox.io/kube-linter/pkg/command/common"
)

func TestAllFormatsSupported(t *testing.T) {
	supportedFormats := set.NewStringSet()
	for format := range formatters {
		supportedFormats.Add(format)
	}
	// TODO: refactor
	assert.ElementsMatch(t, supportedFormats.AsSlice(), []string{common.PlainFormat, common.MarkdownFormat})
}
