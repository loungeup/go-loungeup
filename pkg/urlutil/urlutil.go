package urlutil

import "net/url"

// SafeJoinPath joins the provided path elements to the base URL.
// If the base URL is invalid, it ignores the base and joins the elements as a relative path.
// This function is similar to url.JoinPath but it ignores the base URL if it is invalid.
//
// Reference: https://cs.opensource.google/go/go/+/refs/tags/go1.24.0:src/net/url/url.go;l=1302
func SafeJoinPath(base string, elements ...string) string {
	parsed, err := url.Parse(base)
	if err != nil {
		return new(url.URL).JoinPath(elements...).String()
	}

	return parsed.JoinPath(elements...).String()
}
