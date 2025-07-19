package objectkinds_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
)

func TestGetPersistentVolumeClaimAPIVersion(t *testing.T) {
	apiVersion := objectkinds.GetPersistentVolumeClaimAPIVersion()
	assert.NotEmpty(t, apiVersion)
}
