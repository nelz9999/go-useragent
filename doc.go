// Package gua is Go User Agent helper
package gua

const (
	uaHeader = "User-Agent"

	// From https://golang.org/src/net/http/request.go|defaultUserAgent
	uaDefaultToken   = "Go-http-client"
	uaDefaultVersion = "1.1"
)
