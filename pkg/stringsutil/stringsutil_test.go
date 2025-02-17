package stringsutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToTitleOrUpper(t *testing.T) {
	tests := map[string]struct {
		in, want string
	}{
		"more uppercase characters":                            {in: "HeLLoWorLD", want: "HELLOWORLD"},
		"more lowercase characters":                            {in: "helloWORLd", want: "Helloworld"},
		"equal uppercase and lowercase characters":             {in: "Li Zo", want: "LI ZO"},
		"empty string":                                         {in: "", want: ""},
		"string with only uppercase characters":                {in: "HELLO", want: "HELLO"},
		"string with only lowercase characters":                {in: "hello", want: "Hello"},
		"string with mixed case and non-alphabetic characters": {in: "HeLLo-wOrLd 123", want: "HELLO-WORLD 123"},
		"string with no alphabetic characters":                 {in: "123-456", want: "123-456"},
		"hangul string":                                        {in: "안녕하세요", want: "안녕하세요"},
	}

	for test, tt := range tests {
		t.Run(test, func(t *testing.T) {
			assert.Equal(t, tt.want, ToTitleOrUpper(tt.in))
		})
	}
}
