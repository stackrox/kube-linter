package util

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestConstructRequiredMapMatcher(t *testing.T) {
	object := lintcontext.Object{
		K8sObject: &v1.Pod{
			ObjectMeta: metaV1.ObjectMeta{
				Annotations: map[string]string{
					"a": "1",
					"c": "2",
					"x": "y",
				},
				Labels: map[string]string{
					"a": "3",
					"c": "4",
					"e": "f",
				},
			},
		},
	}
	tests := []struct {
		key       string
		value     string
		fieldType string
		expected  []diagnostic.Diagnostic
	}{{
		key: "a", value: "1", fieldType: "annotation",
	}, {
		key: "a", value: "3", fieldType: "label",
	}, {
		key: "e", value: "f", fieldType: "annotation",
		expected: []diagnostic.Diagnostic{{Message: `no annotation matching "e=f" found`}},
	}, {
		key: "x", value: "y", fieldType: "label",
		expected: []diagnostic.Diagnostic{{Message: `no label matching "x=y" found`}},
	}, {
		key: "a", value: "", fieldType: "label",
	}, {
		key: "a", value: ".*", fieldType: "label",
	}, {
		key: "a", value: "[0-2]", fieldType: "annotation",
	}, {
		key: "a", value: "[0-2]", fieldType: "label",
		expected: []diagnostic.Diagnostic{{Message: `no label matching "a=[0-2]" found`}},
	}, {
		key: "a", value: "!2", fieldType: "label",
	}, {
		key: "!x", value: "", fieldType: "label",
	}, {
		key: "!x", value: "", fieldType: "annotation",
	}}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s=%s %s", tt.key, tt.value, tt.fieldType), func(t *testing.T) {
			tt := tt
			matcher, err := ConstructRequiredMapMatcher(tt.key, tt.value, tt.fieldType)
			assert.NoError(t, err)
			got := matcher(nil, object)
			assert.Equal(t, tt.expected, got)
		})
	}
}
