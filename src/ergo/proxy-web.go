package main

import (
	"net/http"
	"net/http/httputil"

	log "github.com/sirupsen/logrus"
)

func NewWebProxy(rules *Rules) http.Handler {
	return &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			target, host := rules.Find(req)

			// TODO: add SSL-to-noSSL support (SSL termination)
			// req.URL.Scheme = turl.Scheme

			req.URL.Host = target
			req.Host = host

			log.WithField("url", req.URL.String()).Debug("overwritten URL")
		},
	}
}
