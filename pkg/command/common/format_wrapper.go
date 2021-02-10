package common

import (
	"encoding/json"
	"io"
)

const (
	PlainFormat    = "plain"
	MarkdownFormat = "markdown"
	JsonFormat     = "json"
)

func FormatJson(data interface{}, out io.Writer) error {
	return json.NewEncoder(out).Encode(data)
}
