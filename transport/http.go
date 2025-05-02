package transport

import "net/http"

//go:generate mockgen -source http.go -destination=./mocks/mock_http.go -package=mocks

// HTTPDoer is the interface used to execute request using the HTTP protocol. It wraps the Do method implemented by the
// built-in http.Client.
type HTTPDoer interface {
	Do(request *http.Request) (*http.Response, error)
}
