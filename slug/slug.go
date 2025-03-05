package slug

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// Well-known patterns used to generate slugs.
const (
	Alpha         = `[^a-z]`
	AlphaWithDash = `[^a-z-]`
	AlphaWithDot  = `[^a-z.]`

	Alphanumeric         = `[^a-z0-9]`
	AlphanumericWithDash = `[^a-z0-9-]`
)

func FromValue[T ~string](value T, pattern ...string) T {
	slug, _, err := transform.String(
		transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC),
		string(value),
	)
	if err != nil {
		slug = string(value) // Ignore error.
	}

	slug = strings.ToLower(slug)

	// Use alphanumeric by default.
	if len(pattern) < 1 {
		pattern = []string{Alphanumeric}
	}

	slug = regexp.MustCompile(pattern[0]).ReplaceAllString(slug, "")

	return T(slug)
}

func FromValues[E ~string](values []E, pattern ...string) []E {
	result := []E{}
	for _, value := range values {
		result = append(result, FromValue(value, pattern...))
	}

	return result
}
