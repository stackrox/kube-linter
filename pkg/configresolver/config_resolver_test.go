package configresolver

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
	"golang.stackrox.io/kube-linter/pkg/config"
)

var cfg *config.Config = new(config.Config)

func TestHomeIgnorePath(t *testing.T) {
	cfg.Checks.IgnorePaths = []string{"~/test"}

	home, _ := homedir.Dir()
	paths, err := GetIgnorePaths(cfg)
	assert.NoError(t, err)
	assert.Equal(t, paths[0], home+"/test")
}

func TestHomeGlobIgnorePath(t *testing.T) {
	cfg.Checks.IgnorePaths = []string{"~/*.yaml"}

	home, _ := homedir.Dir()
	paths, err := GetIgnorePaths(cfg)
	assert.NoError(t, err)
	assert.Equal(t, paths[0], home+"/*.yaml")
}

func TestInvalidHomeIgnorePath(t *testing.T) {
	cfg.Checks.IgnorePaths = []string{"~~/test"}

	_, err := GetIgnorePaths(cfg)
	assert.Error(t, err)
}

func TestRelativeIgnorePath(t *testing.T) {
	cfg.Checks.IgnorePaths = []string{"../test"}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	parent := filepath.Dir(wd)
	paths, err := GetIgnorePaths(cfg)
	assert.NoError(t, err)
	assert.Equal(t, paths[0], parent+"/test")
}

func TestRelativeGlobIgnorePath(t *testing.T) {
	cfg.Checks.IgnorePaths = []string{"../*.yaml"}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	parent := filepath.Dir(wd)
	paths, err := GetIgnorePaths(cfg)
	assert.NoError(t, err)
	assert.Equal(t, paths[0], parent+"/*.yaml")
}
