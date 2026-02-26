package lint

import (
	"fmt"
	"io"
	"os"
)

// OutputDestination represents where formatted output should be written
type OutputDestination struct {
	Writer io.WriteCloser
	Path   string // empty for stdout
}

// NewOutputDestination creates an output destination
func NewOutputDestination(path string) (*OutputDestination, error) {
	if path == "" {
		return &OutputDestination{
			Writer: nopWriteCloser{os.Stdout},
			Path:   "",
		}, nil
	}

	file, err := os.Create(path) // #nosec G304 -- User-specified output file path
	if err != nil {
		return nil, fmt.Errorf("failed to create output file %q: %w", path, err)
	}

	return &OutputDestination{Writer: file, Path: path}, nil
}

// Close closes the output destination
func (d *OutputDestination) Close() error {
	return d.Writer.Close()
}

// nopWriteCloser wraps an io.Writer and provides a no-op Close method
type nopWriteCloser struct{ io.Writer }

func (nopWriteCloser) Close() error { return nil }
