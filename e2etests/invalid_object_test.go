//go:build e2e

package e2etests

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Regression test for #591 / #669.
//
// A manifest of a known kind that fails to decode used to be silently downgraded to an
// unstructured object. It then became invisible to every pod-spec based check, and the only
// thing the user saw was an unrelated dangling-service diagnostic about a *different* object.
// KubeLinter must say which file it could not read, without needing --verbose.
func TestKubeLinterReportsUndecodableObjectsWithoutVerbose(t *testing.T) {
	kubeLinterBin := os.Getenv(kubeLinterBinEnv)
	require.NotEmpty(t, kubeLinterBin, "Please set %s", kubeLinterBinEnv)

	cases := map[string]struct {
		manifest    string
		wantInError string
	}{
		// #591: an invalid resource quantity.
		"invalid quantity": {
			manifest: `apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-deployment
spec:
  selector:
    matchLabels:
      name: my-label-value
  template:
    metadata:
      labels:
        name: my-label-value
    spec:
      containers:
        - name: my-container
          image: nginx:latest
          resources:
            limits:
              cpu: 10wrong_unit
`,
			wantInError: "quantities must match",
		},
		// #669: a number where the schema wants a string.
		"number instead of string": {
			manifest: `apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-deployment
spec:
  selector:
    matchLabels:
      name: my-label-value
  template:
    metadata:
      labels:
        name: my-label-value
    spec:
      containers:
        - name: my-container
          image: nginx:latest
          livenessProbe:
            exec:
              command:
                - /usr/bin/echo
                - 1234
`,
			wantInError: "cannot unmarshal number",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "manifest.yaml")
			require.NoError(t, os.WriteFile(path, []byte(c.manifest), 0644))

			// Deliberately no --verbose: the user must be told anyway.
			out, _ := exec.Command(kubeLinterBin, "lint", path).CombinedOutput()
			outAsStr := string(out)

			assert.True(t, strings.Contains(outAsStr, "failed to load object"),
				"expected a load failure to be reported, got: %s", outAsStr)
			assert.True(t, strings.Contains(outAsStr, c.wantInError),
				"expected the report to name the real cause %q, got: %s", c.wantInError, outAsStr)
		})
	}
}
