package stringsutil

import (
	"strings"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Formatter[T ~string] struct{ value T }

func NewFormatter[T ~string](value T) *Formatter[T] {
	return &Formatter[T]{value: T(strings.Clone(string(value)))}
}

func (formatter *Formatter[T]) RemoveSpaces() *Formatter[T] {
	formatter.value = T(strings.ReplaceAll(string(formatter.value), " ", ""))

	return formatter
}

func (formatter *Formatter[T]) ToLower() *Formatter[T] {
	formatter.value = T(strings.ToLower(string(formatter.value)))

	return formatter
}

func (formatter *Formatter[T]) ToUpper() *Formatter[T] {
	formatter.value = T(strings.ToUpper(string(formatter.value)))

	return formatter
}

func (formatter *Formatter[T]) TrimSpaces() *Formatter[T] {
	formatter.value = T(strings.TrimSpace(string(formatter.value)))

	return formatter
}

// Value returns the formatted value.
func (formatter *Formatter[T]) Value() T { return formatter.value }

// ToTitleOrUpper converts the input string to title case or upper case based on its content. If the string contains
// more uppercase characters than lowercase characters, it converts the entire string to uppercase. Otherwise, it
// converts the string to title case.
//
// Example:
//
//	ToTitleOrUpper("hello WORLD") // returns "HELLO WORLD"
//	ToTitleOrUpper("hello world") // returns "Hello World"
func ToTitleOrUpper(value string) string {
	if hasMoreUppercaseThanLowercase(value) {
		return strings.ToUpper(value)
	}

	return cases.Title(language.Und).String(value)
}

// hasMoreUppercaseThanLowercase returns true if the input string contains more uppercase characters than lowercase
// characters.
func hasMoreUppercaseThanLowercase(value string) bool {
	lowerCount, upperCount := 0, 0

	for _, rune := range value {
		if unicode.IsLower(rune) {
			lowerCount++
		} else if unicode.IsUpper(rune) {
			upperCount++
		}
	}

	return upperCount >= lowerCount
}
