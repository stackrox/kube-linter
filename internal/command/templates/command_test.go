package templates

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.stackrox.io/kube-linter/internal/command/common"
	"golang.stackrox.io/kube-linter/internal/set"
)

func TestAllFormatsSupported(t *testing.T) {
	supportedFormats := set.NewStringSet()
	for format := range formatsToRenderFuncs {
		supportedFormats.Add(format)
	}
	assert.ElementsMatch(t, supportedFormats.AsSlice(), common.AllSupportedFormats.AsSlice())
}
