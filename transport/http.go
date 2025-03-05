package transport

import "net/http"

// HTTPDoer is the interface used to execute request using the HTTP protocol. It wraps the Do method implemented by the
// built-in http.Client.
type HTTPDoer interface {
	Do(request *http.Request) (*http.Response, error)
}
