package util

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.stackrox.io/kube-linter/internal/pointers"
)

func TestValueInRange(t *testing.T) {
	for _, testCase := range []struct {
		value         int
		lowerBound    int
		upperBound    *int
		expectedMatch bool
	}{
		{
			value:         100,
			lowerBound:    99,
			expectedMatch: true,
		},
		{
			value:         100,
			lowerBound:    100,
			expectedMatch: true,
		},
		{
			value:      100,
			lowerBound: 101,
		},
		{
			value:         100,
			upperBound:    pointers.Int(101),
			expectedMatch: true,
		},
		{
			value:         100,
			upperBound:    pointers.Int(100),
			expectedMatch: true,
		},
		{
			value:      100,
			upperBound: pointers.Int(99),
		},
		{
			value:         100,
			lowerBound:    100,
			upperBound:    pointers.Int(100),
			expectedMatch: true,
		},
		{
			value:         0,
			upperBound:    pointers.Int(0),
			expectedMatch: true,
		},
		{
			value:      1,
			upperBound: pointers.Int(0),
		},
		{
			value:         100,
			lowerBound:    99,
			upperBound:    pointers.Int(101),
			expectedMatch: true,
		},
		{
			value:      102,
			lowerBound: 99,
			upperBound: pointers.Int(101),
		},
		{
			value:      98,
			lowerBound: 99,
			upperBound: pointers.Int(101),
		},
	} {
		c := testCase
		t.Run(fmt.Sprintf("%+v", c), func(t *testing.T) {
			assert.Equal(t, c.expectedMatch, ValueInRange(c.value, c.lowerBound, c.upperBound))
		})
	}
}
