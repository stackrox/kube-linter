package lint

import (
	"io"
	"os"
	"path/filepath"
	"testing"
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
			if (err != nil) != tt.wantErr {
				t.Errorf("NewOutputDestination() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return // Expected error, test passed
			}

			if dest == nil {
				t.Error("NewOutputDestination() returned nil destination")
				return
			}

			// Check if it's stdout
			if tt.wantStdout && dest.Path != "" {
				t.Errorf("Expected stdout (empty path), got path: %s", dest.Path)
			}

			if !tt.wantStdout && dest.Path != tt.path {
				t.Errorf("Expected path %s, got %s", tt.path, dest.Path)
			}

			// Test writing
			testData := "test data"
			n, err := io.WriteString(dest.Writer, testData)
			if err != nil {
				t.Errorf("Failed to write to destination: %v", err)
			}
			if n != len(testData) {
				t.Errorf("Expected to write %d bytes, wrote %d", len(testData), n)
			}

			// Clean up
			if err := dest.Close(); err != nil {
				t.Errorf("Failed to close destination: %v", err)
			}

			// For file outputs, verify the file was created and contains data
			if !tt.wantStdout {
				content, err := os.ReadFile(tt.path)
				if err != nil {
					t.Errorf("Failed to read output file: %v", err)
				}
				if string(content) != testData {
					t.Errorf("Expected file content %q, got %q", testData, string(content))
				}
			}
		})
	}
}

func TestOutputDestination_Close(t *testing.T) {
	t.Run("closing stdout destination is no-op", func(t *testing.T) {
		dest, err := NewOutputDestination("")
		if err != nil {
			t.Fatalf("NewOutputDestination() error = %v", err)
		}

		if err := dest.Close(); err != nil {
			t.Errorf("Close() on stdout destination should not error, got: %v", err)
		}

		// Should be able to close multiple times without error
		if err := dest.Close(); err != nil {
			t.Errorf("Second Close() on stdout destination should not error, got: %v", err)
		}
	})

	t.Run("closing file destination closes file", func(t *testing.T) {
		tmpFile := filepath.Join(t.TempDir(), "test-close.txt")
		dest, err := NewOutputDestination(tmpFile)
		if err != nil {
			t.Fatalf("NewOutputDestination() error = %v", err)
		}

		// Write some data
		_, err = io.WriteString(dest.Writer, "test")
		if err != nil {
			t.Fatalf("Failed to write: %v", err)
		}

		// Close the file
		if err := dest.Close(); err != nil {
			t.Errorf("Close() error = %v", err)
		}

		// Verify file exists and is readable
		// #nosec G304 -- Test file path from controlled test directory
		if _, err := os.ReadFile(tmpFile); err != nil {
			t.Errorf("Failed to read file after close: %v", err)
		}
	})
}
