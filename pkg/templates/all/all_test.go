package all

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.stackrox.io/kube-linter/pkg/templates"
)

func TestTemplatesAreValid(t *testing.T) {
	for _, template := range templates.List() {
		t.Run(template.HumanName, func(t *testing.T) {
			assert.NotEmpty(t, template.HumanName, "human name")
			assert.NotEmpty(t, template.Key, "name")
			assert.NotEmpty(t, template.Description, "description")
			assert.NotNil(t, template.ParseAndValidateParams, "parse and validate params")
			assert.NotNil(t, template.Parameters, "params") // We want people to use the generated code and explicitly set it to an empty list.
			assert.NotNil(t, template.Instantiate, "instantiate")
		})
	}
}
