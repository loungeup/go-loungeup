package maybe

import (
	"time"

	"github.com/google/uuid"
	"github.com/loungeup/go-loungeup/pkg/translations"
)

func Int(value, defaultValue int) int {
	return getValueOrDefault(value, defaultValue, func(value int) bool { return value != 0 })
}

func String(value, defaultValue string) string {
	return getValueOrDefault(value, defaultValue, func(value string) bool { return value != "" })
}

func Time(value, defaultValue time.Time) time.Time {
	return getValueOrDefault(value, defaultValue, func(v time.Time) bool { return !v.IsZero() })
}

func Translations(value, defaultValue translations.Translations) translations.Translations {
	return getValueOrDefault(value, defaultValue, func(value translations.Translations) bool {
		return value != nil && !value.IsZero()
	})
}

func UUID(value, defaultValue uuid.UUID) uuid.UUID {
	return getValueOrDefault(value, defaultValue, func(value uuid.UUID) bool { return value != uuid.Nil })
}

func getValueOrDefault[T any](value, defaultValue T, f func(value T) bool) T {
	if f(value) {
		return value
	}

	return defaultValue
}
