package slicesutil

func Map[S ~[]E, E, M any](s S, f func(E) M) []M {
	result := make([]M, len(s))
	for i, v := range s {
		result[i] = f(v)
	}

	return result
}

func ToAny[S ~[]E, E any](s S) []any {
	result := make([]any, len(s))
	for i, v := range s {
		result[i] = v
	}

	return result
}
