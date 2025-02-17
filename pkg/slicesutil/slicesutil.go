package slicesutil

import "slices"

// Find an element in a slice based on a callback function.
func Find[S ~[]E, E any](s S, findFunc func(e E) bool) (E, bool) {
	if index := slices.IndexFunc(s, findFunc); index != -1 {
		return s[index], true
	}

	var empty E

	return empty, false
}

// Filter a slice based on a callback function.
func Filter[S ~[]E, E any](s S, filterFunc func(element E) bool) []E {
	result := make([]E, 0)

	for _, element := range s {
		if filterFunc(element) {
			result = append(result, element)
		}
	}

	return result
}

// Intersect returns the intersection of two slices of comparable elements.
func Intersect[S ~[]E, E comparable](a, b S) S {
	return IntersectFunc(a, b, func(a, b E) bool { return a == b })
}

// IntersectFunc returns the intersection of two slices using a comparison function.
func IntersectFunc[S ~[]E, E any](a, b S, compareFunc func(E, E) bool) S {
	result := make(S, 0, min(len(a), len(b)))
	if a == nil || b == nil {
		return result
	}

	for _, e := range a {
		if slices.IndexFunc(b, func(be E) bool { return compareFunc(be, e) }) != -1 {
			result = append(result, e)
		}
	}

	return result
}

// Map the slice with the given function.
func Map[S ~[]E, E, M any](s S, mapFunc func(E) M) []M {
	result := make([]M, len(s))
	for i, v := range s {
		result[i] = mapFunc(v)
	}

	return result
}

// Merge two slices of comparable elements using a merge function.
func Merge[S ~[]E, E comparable](a, b S, mergeFunc func(a, b E) E) S {
	return MergeFunc(a, b, func(a, b E) bool { return a == b }, mergeFunc)
}

// MergeFunc merges two slices using a comparison function and a merge function.
func MergeFunc[S ~[]E, E any](a, b S, compareFunc func(a, b E) bool, mergeFunc func(a, b E) E) S {
	combined := slices.Concat(a, b)
	if len(combined) == 0 {
		return nil
	}

	result := make(S, 0, len(combined))
	for _, element := range combined {
		if existingElementIndex := slices.IndexFunc(result, func(existing E) bool {
			return compareFunc(existing, element)
		}); existingElementIndex == -1 {
			result = append(result, element)
		}
	}

	return result
}

// Partition a slice into slices of the given size.
func Partition[S ~[]E, E any](s S, size int) []S {
	if len(s) == 0 || size <= 0 {
		return []S{}
	}

	result := make([]S, 0)

	for i := 0; i < len(s); i += size {
		end := i + size
		if end > len(s) {
			end = len(s)
		}

		result = append(result, s[i:end])
	}

	return result
}

// RemoveEmpty elements from a slice.
func RemoveEmpty[S ~[]E, E comparable](s S) S {
	var empty E

	result := make(S, 0, len(s))

	for _, e := range s {
		if e == empty {
			continue
		}

		result = append(result, e)
	}

	return result
}

func Split[S ~[]E, E any](s S, split func(E) bool) (matched, exclused S) {
	matched, excluded := S{}, S{}

	for _, e := range s {
		if split(e) {
			matched = append(matched, e)
		} else {
			excluded = append(excluded, e)
		}
	}

	return matched, excluded
}

// Subtract returns the difference between two slices of comparable elements. We get the elements from a that are not in
// b.
func Subtract[S ~[]E, E comparable](a, b S) S {
	return SubtractFunc(a, b, func(a, b E) bool { return a == b })
}

// SubtractFunc returns the difference between two slices using a comparison function. We get the elements from a that
// are not in b.
func SubtractFunc[S ~[]E, E any](a, b S, compareFunc func(E, E) bool) S {
	if len(a) == 0 {
		return make(S, 0)
	}

	if len(b) == 0 {
		return a
	}

	result := make(S, 0, len(a))

	for _, aElement := range a {
		if !slices.ContainsFunc(b, func(bElement E) bool { return compareFunc(aElement, bElement) }) {
			result = append(result, aElement)
		}
	}

	return result
}

// ToAny converts a slice to a slice of type []any.
func ToAny[S ~[]E, E any](s S) []any {
	result := make([]any, len(s))
	for i, v := range s {
		result[i] = v
	}

	return result
}
