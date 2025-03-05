package slug

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromValue(t *testing.T) {
	testCases := map[string]struct {
		in      string
		pattern string
		want    string
	}{
		"simple":                  {in: "Foo Bar", want: "foobar"},
		"with dash":               {in: "foo-bar", pattern: AlphaWithDash, want: "foo-bar"},
		"with dashes":             {in: "Foo-Bar", pattern: Alphanumeric, want: "foobar"},
		"with dot":                {in: "foo.bar", pattern: AlphaWithDot, want: "foo.bar"},
		"with special characters": {in: "Fôö Bâr", pattern: Alphanumeric, want: "foobar"},
		"with underscores":        {in: "Foo_Bar", pattern: Alphanumeric, want: "foobar"},
		"with whitespaces":        {in: "F o  oB  a r", pattern: Alphanumeric, want: "foobar"},
	}

	for test, tt := range testCases {
		t.Run(test, func(t *testing.T) {
			if tt.pattern == "" {
				assert.Equal(t, tt.want, FromValue(tt.in))
			} else {
				assert.Equal(t, tt.want, FromValue(tt.in, tt.pattern))
			}
		})
	}
}
