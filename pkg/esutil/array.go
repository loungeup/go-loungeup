package esutil

import (
	"bytes"
	"encoding/json"
)

// Array represents an Elasticsearch array. In Elasticsearch, there is no dedicated array data type. Any field can
// contain zero or more values by default, however, all values in the array must be of the same data type. This type
// automatically handles the case when the array contains only one element.
//
// Reference: https://www.elastic.co/guide/en/elasticsearch/reference/current/array.html
type Array[T any] []T

var _ (json.Marshaler) = (*Array[any])(nil)

func (a Array[T]) MarshalJSON() ([]byte, error) {
	if len(a) == 1 {
		return json.Marshal(a[0])
	}

	return json.Marshal([]T(a))
}

var _ (json.Unmarshaler) = (*Array[any])(nil)

func (a *Array[T]) UnmarshalJSON(data []byte) error {
	if isArrayData(data) {
		return json.Unmarshal(data, (*[]T)(a))
	}

	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	*a = []T{value}

	return nil
}

func isArrayData(data []byte) bool { return bytes.HasPrefix(data, []byte("[")) }
