package timeutil

import (
	"encoding/json"
	"time"
)

// SafeParse is a wrapper around time.Parse that returns a zero time.Time if the parsing fails.
func SafeParse(layout, value string) time.Time {
	result, _ := time.Parse(layout, value)
	return result
}

func Format(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.Format(time.RFC3339)
}

func MostRecent(times ...time.Time) time.Time {
	result := time.Time{}

	for _, t := range times {
		if t.After(result) {
			result = t
		}
	}

	return result
}

func Oldest(times ...time.Time) time.Time {
	result := time.Time{}

	for _, t := range times {
		if result.IsZero() || t.Before(result) {
			result = t
		}
	}

	return result
}

type Date time.Time

func NewDate(t time.Time) Date { return Date(t) }

func (d *Date) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		*d = Date{}
		return nil
	}

	var dateAsString string
	if err := json.Unmarshal(data, &dateAsString); err != nil {
		return err
	}

	if dateAsString == "" {
		*d = Date{}
		return nil
	}

	parsedDate, err := time.Parse(time.DateOnly, dateAsString)
	if err != nil {
		return err
	}

	*d = Date(parsedDate)

	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(d).Format(time.DateOnly))
}
