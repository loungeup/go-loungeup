package stringsutil

import "strings"

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

func (formatter *Formatter[T]) TrimSpaces() *Formatter[T] {
	formatter.value = T(strings.TrimSpace(string(formatter.value)))
	return formatter
}

func (formatter *Formatter[T]) Value() T { return formatter.value }
