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
		Path         []string
		Expected     string
		ErrorExpeted bool
	}{
		{[]string{"~/test"}, home + "/test", false},
		{[]string{"~/*.yaml"}, home + "/*.yaml", false},
		{[]string{"~~/test"}, "", true},
		{[]string{"../test"}, parent + "/test", false},
		{[]string{"../*.yaml"}, parent + "/*.yaml", false},
	}

	for _, e := range tests {
		c.Checks.IgnorePaths = e.Path
		paths, err := GetIgnorePaths(c)

		if e.ErrorExpeted {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, paths[0], e.Expected)
		}
	}
}
