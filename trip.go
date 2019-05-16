package gua

import "net/http"

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (fn roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return fn(r)
}

// NewUARoundTripper will use the given string as the User-Agent header on
// all requests that don't have that header set already.
func NewUARoundTripper(next http.RoundTripper, ua string) http.RoundTripper {
	return roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		if r.Header.Get(uaHeader) == "" {
			r.Header.Add(uaHeader, ua)
		}
		return next.RoundTrip(r)
	})
}
