package client

// RESDataValue is a generic type for a RES data value (https://resgate.io/docs/specification/res-protocol/#data-values).
type RESDataValue[T any] struct {
	Data T `json:"data"`
}

func NewRESDataValue[T any](d T) RESDataValue[T] { return RESDataValue[T]{Data: d} }
