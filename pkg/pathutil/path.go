package pathutil

import (
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
)

// GetAbsolutPath returns the absolute representation of given path.
func GetAbsolutPath(path string) (string, error) {
	switch {
	case path[0] == '~':
		expandedPath, err := homedir.Expand(path)
		if err != nil {
			return "", errors.Wrapf(err, "could not expand path: %q", expandedPath)
		}
		return expandedPath, nil
	case !filepath.IsAbs(path):
		absPath, err := filepath.Abs(path)
		if err != nil {
			return "", errors.Wrapf(err, "could not expand non-absolute path: %q", absPath)
		}
		return absPath, nil
	default:
		return path, nil
	}
}
