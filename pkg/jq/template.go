package jq

import (
	"bytes"
	"encoding/json"
	"io"
)

// Template is a JSON template.
type Template json.RawMessage

func NewTemplate(template json.RawMessage) Template { return Template(template) }

// Reader returns a reader for the template.
func (t Template) Reader() io.Reader { return bytes.NewReader(t) }

// Replace all occurrences of a key with a value in the template.
func (t Template) Replace(key string, value json.RawMessage) Template {
	result := bytes.Clone(t)
	result = bytes.ReplaceAll(result, json.RawMessage("{{"+key+"}}"), value)
	result = bytes.ReplaceAll(result, json.RawMessage("{{ "+key+" }}"), value)

	return result
}
