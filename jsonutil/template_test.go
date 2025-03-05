package jsonutil

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTemplateReplace(t *testing.T) {
	template := NewTemplate(json.RawMessage(`{
		"entityId": "{{entityId}}",
		"entityType": "{{ entityType }}",
	}`)).
		Replace("entityId", json.RawMessage(`123`)).
		Replace("entityType", json.RawMessage(`account`))

	want := `{
		"entityId": "123",
		"entityType": "account",
	}`

	assert.Equal(t, want, template.String())
}
