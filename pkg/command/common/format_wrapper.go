package common

import (
	"encoding/json"
	"io"

	"github.com/pkg/errors"
)

// FormatType defines possible output formats.
type FormatType string

const (
	// PlainFormat is for plain terminal output, possibly colored.
	PlainFormat = "plain"
	// MarkdownFormat is for markdown output suitable for `*.md` files.
	MarkdownFormat = "markdown"
	// JSONFormat is for JSON output.
	JSONFormat = "json"
)

// FormatFunc sets contract formatter of each FormatType should follow.
type FormatFunc func(io.Writer, interface{}) error

// Formatters struct provides way to define supported output formats for a command.
type Formatters struct {
	Formatters map[FormatType]FormatFunc
}

// GetEnabledFormatters returns a string slice enumerating all enabled formatters of this instance.
func (f Formatters) GetEnabledFormatters() []string {
	keys := make([]string, 0, len(f.Formatters))
	for k := range f.Formatters {
		keys = append(keys, string(k))
	}
	return keys
}

// FormatterByType looks up formatter for a given type among the ones configured in this instance.
func (f Formatters) FormatterByType(t string) (FormatFunc, error) {
	formatter := f.Formatters[FormatType(t)]
	if formatter == nil {
		return nil, errors.Errorf("unknown format: %q", t)
	}
	return formatter, nil
}

// FormatJSON formats data as JSON, i.e. implements JSONFormat.
func FormatJSON(out io.Writer, data interface{}) error {
	return json.NewEncoder(out).Encode(data)
}

// Verify that FormatJSON follows the contract of FormatFunc.
var _ FormatFunc = FormatJSON
