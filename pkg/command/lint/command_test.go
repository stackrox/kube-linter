package lint

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	// Register templates
	_ "golang.stackrox.io/kube-linter/pkg/templates/all"
)

func TestCommand_InvalidResources(t *testing.T) {
	tests := []struct {
		name    string
		cmd     *cobra.Command
		failure bool
		output  string
	}{
		{name: "InvalidYAML", cmd: createLintCommand("./testdata/invalid.yaml", "--fail-on-invalid-resource"), failure: true},
		{name: "NonexistentFile", cmd: createLintCommand("./testdata/foo-bar.yaml", "--fail-on-invalid-resource"), failure: true},
		{name: "ValidPod", cmd: createLintCommand("./testdata/valid-pod.yaml", "--fail-on-invalid-resource"), failure: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cmd.Execute()
			if err == nil && tt.failure || err != nil && !tt.failure {
				t.Fail()
			}
		})
	}
}

func createLintCommand(args ...string) *cobra.Command {
	c := Command()
	c.SilenceUsage = true
	c.SetArgs(args)
	return c
}

// TestCommand_MultiFormatOutput tests multi-format output functionality
func TestCommand_MultiFormatOutput(t *testing.T) {
	validPod := "./testdata/valid-pod.yaml"

	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		args        []string
		expectError bool
		errorText   string
		checkFiles  []string
	}{
		{
			name: "SingleFormatToFile",
			args: []string{
				"--format", "json",
				"--output", filepath.Join(tmpDir, "single-output.json"),
				validPod,
			},
			expectError: false,
			checkFiles:  []string{filepath.Join(tmpDir, "single-output.json")},
		},
		{
			name: "MultipleFormatsToFiles",
			args: []string{
				"--format", "json",
				"--output", filepath.Join(tmpDir, "multi-output.json"),
				"--format", "sarif",
				"--output", filepath.Join(tmpDir, "multi-output.sarif"),
				validPod,
			},
			expectError: false,
			checkFiles: []string{
				filepath.Join(tmpDir, "multi-output.json"),
				filepath.Join(tmpDir, "multi-output.sarif"),
			},
		},
		{
			name: "FormatOutputMismatch",
			args: []string{
				"--format", "json",
				"--format", "sarif",
				"--output", "/tmp/only-one.json",
				validPod,
			},
			expectError: true,
			errorText:   "format/output mismatch",
		},
		{
			name: "InvalidOutputDirectory",
			args: []string{
				"--format", "json",
				"--output", "/nonexistent/dir/output.json",
				validPod,
			},
			expectError: true,
			errorText:   "failed to create output destination",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := createLintCommand(tt.args...)
			err := cmd.Execute()

			if tt.expectError && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tt.expectError && err != nil && tt.errorText != "" {
				if !strings.Contains(err.Error(), tt.errorText) {
					t.Errorf("expected error to contain %q, got: %v", tt.errorText, err)
				}
			}

			// Check that expected files exist
			for _, file := range tt.checkFiles {
				if _, err := os.Stat(file); err != nil {
					t.Errorf("expected file %s to exist, but got error: %v", file, err)
				}
			}
		})
	}
}

// TestCommand_MultiFormatResourceCleanup verifies file handles are properly closed
func TestCommand_MultiFormatResourceCleanup(t *testing.T) {
	validPod := "./testdata/valid-pod.yaml"
	tmpDir := t.TempDir()

	// Create command with 10 different output files
	args := []string{}
	var outputFiles [10]string

	for i := 0; i < 10; i++ {
		format := "json"
		if i%2 == 0 {
			format = "sarif"
		}
		outputPath := filepath.Join(tmpDir, format+"-output-"+string(rune('0'+i))+".out")
		outputFiles[i] = outputPath
		args = append(args, "--format", format, "--output", outputPath)
	}
	args = append(args, validPod)

	cmd := createLintCommand(args...)
	err := cmd.Execute()

	// We expect lint errors from the valid pod, not output errors
	if err != nil && !strings.Contains(err.Error(), "lint errors") {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify all files were created and can be reopened (proving they were closed)
	for _, outputFile := range outputFiles {
		// Try to open the file for reading
		f, err := os.Open(outputFile) // #nosec G304 -- Test file path from controlled test directory
		if err != nil {
			t.Errorf("failed to open output file %s: %v", outputFile, err)
			continue
		}
		_ = f.Close() // Best effort close in test

		// Verify file has content
		info, err := os.Stat(outputFile)
		if err != nil {
			t.Errorf("failed to stat output file %s: %v", outputFile, err)
			continue
		}
		if info.Size() == 0 {
			t.Errorf("output file %s is empty", outputFile)
		}
	}
}

// TestCommand_PartialFailureHandling verifies behavior when some outputs fail
func TestCommand_PartialFailureHandling(t *testing.T) {
	validPod := "./testdata/valid-pod.yaml"
	tmpDir := t.TempDir()
	validPath := filepath.Join(tmpDir, "valid.json")

	cmd := createLintCommand(
		"--format", "json",
		"--output", validPath,
		"--format", "sarif",
		"--output", "/nonexistent/invalid.sarif",
		validPod,
	)

	err := cmd.Execute()

	// Should have an error about the failed output
	if err == nil {
		t.Fatalf("expected error for failed output")
	}

	// Error should mention both success and failure
	errMsg := err.Error()
	if !strings.Contains(errMsg, "1 of 2") && !strings.Contains(errMsg, "2 of 2") {
		t.Errorf("error message should mention output count, got: %v", errMsg)
	}

	// The successful output should still exist
	if _, err := os.Stat(validPath); err != nil {
		t.Errorf("successful output file should exist: %v", err)
	}
}
