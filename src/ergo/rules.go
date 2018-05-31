package main

import (
	"bufio"
	"io"
	"net/http"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Rules struct {
	rs []*regexp.Regexp // rules
	ts []string         // targets
	hs []string         // hostnames
}

func NewRules(file io.Reader) *Rules {
	rules := &Rules{
		rs: make([]*regexp.Regexp, 0),
		ts: make([]string, 0),
		hs: make([]string, 0),
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.Split(line, "-->")
		if len(parts) < 2 {
			log.WithField("line", line).Fatal("wrong format")
		}

		part := strings.TrimSpace(parts[0])
		re, err := regexp.Compile(part)
		if err != nil {
			log.WithError(err).WithField("rule", part).Fatal("can't compile regexp")
		}

		rules.rs = append(rules.rs, re)

		parts = strings.Split(parts[1], " as ")
		target := strings.TrimSpace(parts[0])
		host := ""
		if len(parts) > 1 {
			host = strings.TrimSpace(parts[1])
		}

		rules.ts = append(rules.ts, target)
		rules.hs = append(rules.hs, host)
	}

	return rules
}

// Find returns overwrited hostname for given request
func (r *Rules) Find(req *http.Request) (string, string) {
	url := r.GetAbsURL(req)

	for n, rule := range r.rs {
		if rule.FindString(url) != "" {
			target := r.ts[n]

			host := r.hs[n]
			if host == "" {
				host = req.Host
			}

			log.WithField("url", url).WithField("target", r.ts[n]).WithField("host", host).Debug("overwritten target and host for url")

			return target, host
		}
	}

	log.WithField("url", url).WithField("host", req.Host).Debug("no rules matched")
	return req.Host, req.Host // nothing found, return original hostname
}

// GetRequestURI returns absolute full URL for request
// NOTE: we can't use req.URL.String() because sometimes there is no host in req.URL
func (r *Rules) GetAbsURL(req *http.Request) string {
	url := req.Host + req.URL.Path

	if req.URL.RawQuery != "" || req.URL.ForceQuery {
		url += "?" + req.URL.RawQuery
	}

	if req.URL.Scheme != "" {
		url = req.URL.Scheme + "://" + url
	}

	return url
}
