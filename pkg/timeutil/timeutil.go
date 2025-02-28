package timeutil

import (
	"encoding/json"
	"time"
)

func FormatIfNotZero(layout string, t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.Format(layout)
}

// SafeParse is a wrapper around time.Parse that returns a zero time.Time if the parsing fails.
func SafeParse(layout, value string) time.Time {
	result, _ := time.Parse(layout, value)

	return result
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

type RFC3339Time time.Time

func NewRFC3339Time(t time.Time) RFC3339Time { return RFC3339Time(t) }

func NowRFC3339() RFC3339Time { return NewRFC3339Time(time.Now()) }

func (t RFC3339Time) Format() string { return t.Time().Format(time.RFC3339) }

func (t RFC3339Time) Time() time.Time { return time.Time(t) }

func (t RFC3339Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Format())
}

func (t *RFC3339Time) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		*t = RFC3339Time{}

		return nil
	}

	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return err
	}

	*t = NewRFC3339Time(parsed)

	return nil
}

type DateOnlyTime time.Time

func NewDateOnlyTime(t time.Time) DateOnlyTime { return DateOnlyTime(t) }

func (t DateOnlyTime) Format() string { return t.Time().Format(time.DateOnly) }

func (t DateOnlyTime) Time() time.Time { return time.Time(t) }

func (t DateOnlyTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Format())
}

func (t *DateOnlyTime) UnmarshalJSON(text []byte) error {
	if len(text) == 0 {
		*t = DateOnlyTime{}

		return nil
	}

	var value string
	if err := json.Unmarshal(text, &value); err != nil {
		return err
	}

	parsed, err := time.Parse(time.DateOnly, value)
	if err != nil {
		return err
	}

	*t = NewDateOnlyTime(parsed)

	return nil
}
