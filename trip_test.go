package gua

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRoundTripper(t *testing.T) {
	testCases := map[string]struct {
		client *http.Client
		decor  func(*http.Request) *http.Request
		expect string
	}{
		"undefined": {
			http.DefaultClient,
			func(r *http.Request) *http.Request { return r },
			uaDefaultToken,
		},
		"fallback": {
			&http.Client{Transport: NewUARoundTripper(http.DefaultTransport, "Client/0.1")},
			func(r *http.Request) *http.Request { return r },
			"Client",
		},
		"override": {
			&http.Client{Transport: NewUARoundTripper(http.DefaultTransport, "Client/0.1")},
			func(r *http.Request) *http.Request {
				r.Header.Add("user-agent", "Per-Request/10.9.8")
				return r
			},
			"Per-Request",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var ua string
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ua = r.Header.Get("User-Agent")
			}))
			defer srv.Close()

			req, err := http.NewRequest("GET", srv.URL, nil)
			if err != nil {
				t.Fatalf("%v\n", err)
			}

			_, err = tc.client.Do(tc.decor(req))
			if err != nil {
				t.Fatalf("%v\n", err)
			}

			if !strings.Contains(ua, tc.expect) {
				t.Errorf("expected %q in %q\n", tc.expect, ua)
			}
			// t.Logf("UA: %s\n", ua)
		})
	}
}
