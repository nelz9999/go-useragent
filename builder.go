package gua

import (
	"fmt"
	"strings"
)

const (
	disallowed           = `"(),/:;<=>?@[\]{}`
	msgTokenWhitespace   = "token contains whitespace"
	msgTokenInvalid      = "token contains invalid character"
	msgRequireToken      = "version requires token"
	msgVersionWhitespace = "version contains whitespace"
	msgVersionInvalid    = "version contains invalid character"
)

type product struct {
	token    string
	version  string
	comments []string
}

// check enforces product/product-version/token characters
// as described in this blog post:
// http://www.bizcoder.com/the-much-maligned-user-agent-header
// This is not intended to be a FULL RFC check, but should be
// mostly functional.
func (p *product) check() error {
	if len(strings.Fields(p.token)) > 1 {
		return fmt.Errorf("%s: %q", msgTokenWhitespace, p.token)
	}
	if strings.ContainsAny(p.token, disallowed) {
		return fmt.Errorf("%s: %q", msgTokenInvalid, p.token)
	}
	if len(p.token) == 0 && len(p.version) > 0 {
		return fmt.Errorf("%s: %q", msgRequireToken, p.version)
	}
	if len(strings.Fields(p.version)) > 1 {
		return fmt.Errorf("%s: %q", msgVersionWhitespace, p.version)
	}
	if strings.ContainsAny(p.version, disallowed) {
		return fmt.Errorf("%s: %q", msgVersionInvalid, p.version)
	}
	return nil
}

// Builder holds zero to many product stanzas
type Builder struct {
	parts []*product
}

// NewBuilder is the starter for a fluent User-Agent assembler
func NewBuilder() *Builder {
	return &Builder{}
}

// With creates a new product stanza within a User-Agent string.
// Generally the stanza will look like "token/version (comment0; comment1; etc)"
// The return value is itself so you can append more stanza.
func (b *Builder) With(token string, version string, comments ...string) *Builder {
	b.parts = append(b.parts, &product{
		token:    token,
		version:  version,
		comments: comments,
	})
	return b
}

// Default adds the equivalent to the default net/http User-Agent string
func (b *Builder) Default() *Builder {
	return b.With(uaDefaultToken, uaDefaultVersion)
}

// Build assembles the different stanzas into a whole, checking each one for
// relative conformance to the RFC 7231.
func (b *Builder) Build() (string, error) {
	var out strings.Builder
	for _, prd := range b.parts {
		if err := prd.check(); err != nil {
			return "", err
		}
		if len(prd.token) > 0 {
			if out.Len() > 0 {
				out.WriteString(" ")
			}
			out.WriteString(prd.token)
			if len(prd.version) > 0 {
				out.WriteString("/")
				out.WriteString(prd.version)
			}
		}
		if len(prd.comments) > 0 {
			cmts := strings.Join(prd.comments, "; ")
			if len(cmts) > 0 {
				out.WriteString(" (")
				out.WriteString(cmts)
				out.WriteString(")")
			}
		}
	}
	return out.String(), nil
}
