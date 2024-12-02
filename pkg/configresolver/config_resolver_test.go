package configresolver

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
	"golang.stackrox.io/kube-linter/pkg/config"
)

func TestIgnorePaths(t *testing.T) {
	home, homeErr := homedir.Dir()
	if homeErr != nil {
		t.Fatal(homeErr)
	}
	wd, wdErr := os.Getwd()
	if wdErr != nil {
		t.Fatal(wdErr)
	}

	parent := filepath.Dir(wd)
	c := new(config.Config)

	var tests = []struct {
		Paths        []string
		Expected     string
		ErrorExpeted bool
	}{
		{[]string{"~/test"}, home + "/test", false},
		{[]string{"~/*.yaml"}, home + "/*.yaml", false},
		{[]string{"~~/test"}, "", true},
		{[]string{"../test"}, parent + "/test", false},
		{[]string{"../*.yaml"}, parent + "/*.yaml", false},
		{[]string{"~/test", "~/test"}, home + "/test", false},
	}

	for _, e := range tests {
		c.Checks.IgnorePaths = e.Paths
		paths, err := GetIgnorePaths(c)

		if e.ErrorExpeted {
			assert.Error(t, err)
		} else {
			for _, path := range paths {
				assert.NoError(t, err)
				assert.Equal(t, e.Expected, path)
			}
		}
	}
}
