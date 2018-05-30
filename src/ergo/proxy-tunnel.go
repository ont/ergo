package main

import (
	"io"
	"net"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func NewTunnelProxy(rules *Rules, isConnect bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var host string

		if isConnect {
			// TODO: better solution with w.WriteHeader(http.StatusOK) + internal roundrip handling after hijack
			host = "localhost:" + ERGO_PORT // connect to self for further upstream selection by rules
		} else {
			host = rules.Find(r)
		}

		dst, err := net.DialTimeout("tcp", host, 10*time.Second)
		if err != nil {
			http.Error(w, "Error contacting backend server.", 500)
			log.WithError(err).WithField("target", host).Error("Error dialing tunnel backend")
			return
		}

		// after successful connection to destination we must inform client about successful state of CONNECT command
		if isConnect {
			w.WriteHeader(http.StatusOK)
		}

		hj, ok := w.(http.Hijacker)
		if !ok {
			http.Error(w, "hijack conversion error", 500)
			log.WithError(err).Error("hijack conversion error")
			return
		}

		cli, _, err := hj.Hijack()
		if err != nil {
			http.Error(w, "hijack error", 500)
			log.WithError(err).Error("hijack error")
			return
		}
		defer cli.Close()
		defer dst.Close()

		// for non-CONNECT methods we resend original request from client to destination
		if !isConnect {
			err = r.Write(dst)
			if err != nil {
				log.WithError(err).Error("error copying request to target")
				return
			}
		}

		errc := make(chan error, 2)
		cp := func(dst io.Writer, src io.Reader) {
			_, err := io.Copy(dst, src)
			errc <- err
		}
		go cp(dst, cli)
		go cp(cli, dst)
		<-errc
	})
}
