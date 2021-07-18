package main

import (
	"log"
	"net/http"
	"os"

	logging "github.com/op/go-logging"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Log - use same logger throughout the project
	Log          = logging.MustGetLogger("logger")
	format       = logging.MustStringFormatter(`%{color}%{time:15:04:05.000} %{level} %{shortfunc}%{color:reset} - %{message}`)
	backend1     = logging.NewLogBackend(os.Stdout, "", 0)
	formattedLog = logging.NewBackendFormatter(backend1, format)
)

func collectMetrics() {
	for _, t := range cfg.Targets {
		t.getMetrics()
	}
}

func httpWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		collectMetrics()
		h.ServeHTTP(w, r)
	})
}

func main() {
	logging.SetBackend(formattedLog)
	logging.SetLevel(logging.NOTICE, "logger")
	getConfig()
	logLevel := getLogLevel(cfg.Server.LogLevel)
	logging.SetLevel(logLevel, "logger")
	Log.Info("Started Main")

	collectMetrics()

	http.Handle(getHTTPEndpoint(), httpWrapper(promhttp.Handler()))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html>
			<head><title>Graphite-Exporter</title></head>
			<body>
			<h1>Graphite Exporter</h1>
			<p><a href="` + getHTTPEndpoint() + `">Metrics</a></p>
			</body>
			</html>`))
	})
	Log.Info("Starting http server")
	log.Fatal(http.ListenAndServe(getHTTPPort(), nil))
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
