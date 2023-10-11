package pointer

func From[T any](v T) *T { return &v }
