package common

import (
	"encoding/json"
	"io"

	"github.com/pkg/errors"
)

type FormatType string

const (
	PlainFormat    = "plain"
	MarkdownFormat = "markdown"
	JsonFormat     = "json"
)

type FormatFunc func(io.Writer, interface{}) error

type Formatters struct {
	Formatters map[FormatType]FormatFunc
}

func (f Formatters) GetEnabledFormatters() []string {
	keys := make([]string, 0, len(f.Formatters))
	for k := range f.Formatters {
		keys = append(keys, (string)(k))
	}
	return keys
}

func (f Formatters) FormatterByType(t string) (FormatFunc, error) {
	formatter := f.Formatters[(FormatType)(t)]
	if formatter == nil {
		return nil, errors.Errorf("unknown format: %q", t)
	}
	return formatter, nil
}

func FormatJson(out io.Writer, data interface{}) error {
	return json.NewEncoder(out).Encode(data)
}
