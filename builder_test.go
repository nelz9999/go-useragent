package gua

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"
)

func TestProductCheckErrors(t *testing.T) {
	white := []byte("\n\t ")
	now := time.Now().UnixNano()
	rnd := rand.New(rand.NewSource(now))

	testCases := map[string]struct {
		prd *product
	}{
		msgTokenWhitespace: {
			&product{
				token: fmt.Sprintf(
					"example%c%d",
					white[rnd.Intn(len(white))],
					now,
				),
				version: "9.8.7",
			},
		},
		msgTokenInvalid: {
			&product{
				token: fmt.Sprintf(
					"example%c%d",
					[]byte(disallowed)[rnd.Intn(len(disallowed))],
					now,
				),
				version: "9.8.7",
			},
		},
		msgRequireToken: {
			&product{
				token: "",
				version: fmt.Sprintf(
					"example%c%d",
					[]byte(disallowed)[rnd.Intn(len(disallowed))],
					now,
				),
			},
		},
		msgVersionWhitespace: {
			&product{
				token: "example",
				version: fmt.Sprintf(
					"example%c%d",
					white[rnd.Intn(len(white))],
					now,
				),
			},
		},
		msgVersionInvalid: {
			&product{
				token: "example",
				version: fmt.Sprintf(
					"example%c%d",
					[]byte(disallowed)[rnd.Intn(len(disallowed))],
					now,
				),
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.prd.check()
			if err == nil {
				t.Fatalf("expected error: %v\n", tc.prd)
			}
			if !strings.Contains(err.Error(), name) {
				t.Errorf("expected %q in %q\n", name, err.Error())
			}
		})
	}
}

func TestBuilder(t *testing.T) {
	testCases := map[string]struct {
		b *Builder
		o string
	}{
		"default": {
			NewBuilder().Default(),
			"Go-http-client/1.1",
		},
		"rfc7231 5.5.3": {
			NewBuilder().
				With("CERN-LineMode", "2.15").
				With("libwww", "2.17b3"),
			"CERN-LineMode/2.15 libwww/2.17b3",
		},
		"simulated Chrome": { // http://www.bizcoder.com/the-much-maligned-user-agent-header
			NewBuilder().
				With("Mozilla", "5.0", "Windows NT 6.3", "WOW64").
				With("AppleWebKit", "537.36", "KHTML, like Gecko").
				With("Chrome", "34.0.1847.131").
				With("Safari", "537.36"),
			"Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/34.0.1847.131 Safari/537.36",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			out, err := tc.b.Build()
			if err != nil {
				t.Fatalf("unexpected: %v\n", err)
			}

			if out != tc.o {
				t.Errorf("expected %q; got %q\n", tc.o, out)
			}
			// t.Logf("ua: %q\n", out)
		})
	}
}
