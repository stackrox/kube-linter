package ignore

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestObjectForCheck(t *testing.T) {
	for _, testCase := range []struct {
		annotations  map[string]string
		checkName    string
		shouldIgnore bool
	}{
		{
			annotations: nil,
			checkName:   "some-check",
		},
		{
			annotations: map[string]string{},
			checkName:   "some-check",
		},
		{
			annotations: map[string]string{
				"random-unrelated":                       "blah",
				"ignore-check.kube-linter.io/some-check": "Not applicable",
			},
			checkName:    "some-check",
			shouldIgnore: true,
		},
		{
			annotations: map[string]string{
				"random-unrelated":                       "blah",
				"ignore-check.kube-linter.io/some-check": "Not applicable",
			},
			checkName: "some-check-2",
		},
		{
			annotations: map[string]string{
				"random-unrelated":                       "blah",
				"ignore-check.kube-linter.io/some-check": "Not applicable",
			},
			checkName: "other-check",
		},
		{
			annotations: map[string]string{
				"random-unrelated":                        "blah",
				"ignore-check.kube-linter.io/some-check":  "Not applicable",
				"ignore-check.kube-linter.io/other-check": "Not applicable",
			},
			checkName:    "other-check",
			shouldIgnore: true,
		},
		{
			annotations: map[string]string{
				"random-unrelated":          "blah",
				"kube-linter.io/ignore-all": "Too much of a mess",
			},
			checkName:    "other-check",
			shouldIgnore: true,
		},
		{
			annotations: map[string]string{
				"random-unrelated":          "blah",
				"kube-linter.io/ignore-all": "Too much of a mess",
			},
			checkName:    "some-other-check",
			shouldIgnore: true,
		},
	} {
		c := testCase
		t.Run(fmt.Sprintf("%+v", c), func(t *testing.T) {
			assert.Equal(t, c.shouldIgnore, ObjectForCheck(c.annotations, c.checkName))
		})
	}
}
