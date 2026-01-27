package lint

import (
	"fmt"
	"strings"

	"golang.stackrox.io/kube-linter/internal/set"
	"golang.stackrox.io/kube-linter/pkg/command/common"
)

// FormatOutputPair represents a format and its output destination
type FormatOutputPair struct {
	Format common.FormatType
	Output string // empty string means stdout
}

// ValidateAndPairFormatsOutputs validates format and output flags and pairs them
func ValidateAndPairFormatsOutputs(formats, outputs, allowedFormats []string) ([]FormatOutputPair, error) {
	if len(formats) == 0 {
		return nil, fmt.Errorf("at least one format must be specified")
	}

	// Validate formats
	allowedSet := set.NewFrozenStringSet(allowedFormats...)
	for _, format := range formats {
		if !allowedSet.Contains(format) {
			return nil, fmt.Errorf("invalid format %q: allowed values are: %s",
				format, strings.Join(allowedFormats, ", "))
		}
	}

	// Handle backward compatibility: no outputs means stdout
	if len(outputs) == 0 {
		// Prevent multiple formats to stdout (creates unparseable mixed output)
		if len(formats) > 1 {
			return nil, fmt.Errorf("multiple formats require explicit --output flags. " +
				"Use --output to specify files, or use a single --format for stdout")
		}
		pairs := make([]FormatOutputPair, len(formats))
		for i, format := range formats {
			pairs[i] = FormatOutputPair{
				Format: common.FormatType(format),
				Output: "",
			}
		}
		return pairs, nil
	}

	// Validate output count matches format count
	if len(formats) != len(outputs) {
		return nil, fmt.Errorf("format/output mismatch: %d format(s) specified but %d output(s) provided. "+
			"Each format must have a corresponding output, or omit all outputs to use stdout",
			len(formats), len(outputs))
	}

	// Check for duplicate output files
	outputFiles := make(map[string]bool)
	for _, output := range outputs {
		if output != "" && outputFiles[output] {
			return nil, fmt.Errorf("duplicate output file: %q specified multiple times", output)
		}
		outputFiles[output] = true
	}

	// Pair formats with outputs
	pairs := make([]FormatOutputPair, len(formats))
	for i := range formats {
		pairs[i] = FormatOutputPair{
			Format: common.FormatType(formats[i]),
			Output: outputs[i],
		}
	}

	return pairs, nil
}
