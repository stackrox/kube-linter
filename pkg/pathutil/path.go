package pathutil

import (
	"fmt"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

// GetAbsolutPath returns the absolute representation of given path.
func GetAbsolutPath(path string) (string, error) {
	switch {
	case path[0] == '~':
		expandedPath, err := homedir.Expand(path)
		if err != nil {
			return "", fmt.Errorf("could not expand path: %q: %w", expandedPath, err)
		}
		return expandedPath, nil
	case !filepath.IsAbs(path):
		absPath, err := filepath.Abs(path)
		if err != nil {
			return "", fmt.Errorf("could not expand non-absolute path: %q: %w", absPath, err)
		}
		return absPath, nil
	default:
		return path, nil
	}
}
