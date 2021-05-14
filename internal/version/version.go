package version

import (
	"golang.stackrox.io/kube-linter/pkg/stringutils"
)

var (
	version string //XDef:VERSION
)

// Get returns the version.
func Get() string {
	return stringutils.OrDefault(version, "development")
}
