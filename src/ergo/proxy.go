package main

import (
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Proxy struct {
	webProxy       http.Handler
	connectProxy   http.Handler
	websocketProxy http.Handler
}

func NewProxy(rules *Rules) *Proxy {
	return &Proxy{
		webProxy:       NewWebProxy(rules),
		connectProxy:   NewTunnelProxy(rules, true),
		websocketProxy: NewTunnelProxy(rules, false),
	}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := log.WithField("method", r.Method).WithField("url", r.URL.String())

	switch {
	case r.Method == http.MethodConnect:
		logger.Debug("CONNECT request")
		p.connectProxy.ServeHTTP(w, r)

	case p.isWebsocket(r):
		logger.Debug("websocket request")
		p.websocketProxy.ServeHTTP(w, r)

	default:
		logger.Debug("web request")
		p.webProxy.ServeHTTP(w, r)
	}
}

func (p *Proxy) isWebsocket(req *http.Request) bool {
	res := false
	if strings.ToLower(p.getFirstHeader(req, "Connection")) == "upgrade" {
		res = p.getFirstHeader(req, "Upgrade") == "websocket"
	}

	return res
}

func (p *Proxy) getFirstHeader(req *http.Request, name string) string {
	headers := req.Header[name]
	if len(headers) == 0 {
		return ""
	}
	return headers[0]
}
