package params

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name          string
		params        Params
		expectedError error
	}{
		{
			name: "valid annotation",
			params: Params{
				Annotation: "some-annotation",
			},
			expectedError: nil,
		},
		{
			name: "missing annotation",
			params: Params{
				Annotation: "",
			},
			expectedError: errors.Errorf("invalid parameters: required param annotation not found"),
		},
		{
			name: "valid annotation but with additional invalid param",
			params: Params{
				Annotation: "",
			},
			expectedError: errors.Errorf("invalid parameters: required param annotation not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.params.Validate()
			if tt.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expectedError.Error())
			}
		})
	}
}
