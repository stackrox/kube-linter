package util

import (
	"encoding/json"
	"strings"

	"golang.stackrox.io/kube-linter/pkg/check"
)

// MustParseParameterDesc unmarshals the given JSON into a templates.ParameterDesc.
func MustParseParameterDesc(asJSON string) check.ParameterDesc {
	var out check.ParameterDesc

	decoder := json.NewDecoder(strings.NewReader(asJSON))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&out); err != nil {
		panic(err)
	}
	return out
}
