package models

// DataValue is the same as in the res package but with a generic type.
type DataValue[T any] struct {
	Data T `json:"data"`
}

func NewDataValue[T any](data T) DataValue[T] { return DataValue[T]{Data: data} }
