package util

import (
	"github.com/mitchellh/mapstructure"
)

// DecodeMapStructure decodes the given map[string]interface{} into the given out variable, typically
// a pointer to a struct.
func DecodeMapStructure(m map[string]interface{}, out interface{}) error {
	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		ErrorUnused: true,
		TagName:     "json",
		Result:      out,
	})
	if err != nil {
		return err
	}
	return dec.Decode(m)
}
