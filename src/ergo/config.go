package main

import (
	"crypto/tls"
	"io"
	"net/http"
)

func DisableSSLVerification() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
}

func ConfiguredProxy(file io.Reader) *Proxy {
	rules := NewRules(file)

	return NewProxy(rules)
}
