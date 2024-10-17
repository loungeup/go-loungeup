package timeutil

import "time"

func Format(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.Format(time.RFC3339)
}

func MostRecent(times ...time.Time) time.Time {
	result := time.Time{}

	for _, time := range times {
		if time.After(result) {
			result = time
		}
	}

	return result
}

func Oldest(times ...time.Time) time.Time {
	result := time.Time{}

	for _, time := range times {
		if result.IsZero() || time.Before(result) {
			result = time
		}
	}

	return result
}
