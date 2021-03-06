package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	verbose = kingpin.Flag("verbose", "Verbose logging.").Short('v').Bool()
	config  = kingpin.Flag("config", "Config file with rules").Short('c').Default("/ergo.conf").File()
)

const (
	ERGO_PORT = "2000"
)

func main() {
	kingpin.Parse()

	if *verbose {
		log.SetLevel(log.DebugLevel)
	}

	DisableSSLVerification()

	proxy := ConfiguredProxy(*config)

	log.Fatal(http.ListenAndServe(":"+ERGO_PORT, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})))
}
