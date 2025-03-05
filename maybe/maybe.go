package maybe

import (
	"time"

	"github.com/google/uuid"
	"github.com/loungeup/go-loungeup/translations"
)

type optionalValue[T any] struct {
	valid bool
	value T
}

func (v *optionalValue[T]) Or(value T) T {
	if v.valid {
		return v.value
	}

	return value
}

func Int(value int) *optionalValue[int] {
	return &optionalValue[int]{valid: value != 0, value: value}
}

func String(value string) *optionalValue[string] {
	return &optionalValue[string]{valid: value != "", value: value}
}

func Time(value time.Time) *optionalValue[time.Time] {
	return &optionalValue[time.Time]{valid: !value.IsZero(), value: value}
}

func Translations(value translations.Translations) *optionalValue[translations.Translations] {
	return &optionalValue[translations.Translations]{valid: value != nil && !value.IsZero(), value: value}
}

func UUID(value uuid.UUID) *optionalValue[uuid.UUID] {
	return &optionalValue[uuid.UUID]{valid: value != uuid.Nil, value: value}
}
