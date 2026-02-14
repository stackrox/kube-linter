package lint

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOutputDestination(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		wantStdout  bool
		wantErr     bool
		setupFunc   func(string)
		cleanupFunc func(string)
	}{
		{
			name:       "empty path returns stdout",
			path:       "",
			wantStdout: true,
			wantErr:    false,
		},
		{
			name:       "valid file path",
			path:       filepath.Join(t.TempDir(), "test-output.txt"),
			wantStdout: false,
			wantErr:    false,
		},
		{
			name:       "invalid directory path",
			path:       "/nonexistent/directory/file.txt",
			wantStdout: false,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupFunc != nil {
				tt.setupFunc(tt.path)
			}
			if tt.cleanupFunc != nil {
				defer tt.cleanupFunc(tt.path)
			}

			dest, err := NewOutputDestination(tt.path)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, dest)

			// Check if it's stdout
			if tt.wantStdout {
				assert.Empty(t, dest.Path, "Expected stdout (empty path)")
			} else {
				assert.Equal(t, tt.path, dest.Path)
			}

			// Test writing
			testData := "test data"
			n, err := io.WriteString(dest.Writer, testData)
			require.NoError(t, err)
			assert.Equal(t, len(testData), n)

			// Clean up
			assert.NoError(t, dest.Close())

			// For file outputs, verify the file was created and contains data
			if !tt.wantStdout {
				content, err := os.ReadFile(tt.path)
				require.NoError(t, err)
				assert.Equal(t, testData, string(content))
			}
		})
	}
}

func TestOutputDestination_Close(t *testing.T) {
	t.Run("closing stdout destination is no-op", func(t *testing.T) {
		dest, err := NewOutputDestination("")
		require.NoError(t, err)

		assert.NoError(t, dest.Close())

		// Should be able to close multiple times without error
		assert.NoError(t, dest.Close())
	})

	t.Run("closing file destination closes file", func(t *testing.T) {
		tmpFile := filepath.Join(t.TempDir(), "test-close.txt")
		dest, err := NewOutputDestination(tmpFile)
		require.NoError(t, err)

		// Write some data
		_, err = io.WriteString(dest.Writer, "test")
		require.NoError(t, err)

		// Close the file
		assert.NoError(t, dest.Close())

		// Verify file exists and is readable
		// #nosec G304 -- Test file path from controlled test directory
		_, err = os.ReadFile(tmpFile)
		assert.NoError(t, err)
	})
}
