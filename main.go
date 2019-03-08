package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	config   = Config{}
	graphite = Graphite{}
	gauges   = make(map[string]prometheus.GaugeVec)
	// HTTPPort exposed metrics port
	HTTPPort = 8080
	// HTTPEndpoint exposed metrics endpoint
	HTTPEndpoint = "/metrics"
	namespace    = "graphite_exporter"
	// DebugLogging global tracker for extended logging
	DebugLogging = false
)

func collectMetrics() {
	for _, m := range config.Metrics {
		logMessage("collecting metrics for: %s", m.Name)
		logMessage(" - getting prometheus gauge")
		g := getGauge(m)
		logMessage(" - getting metrics from graphite")
		respSlice := graphite.getMetric(m)
		for _, gr := range respSlice {
			target := trimAndReplace(gr.Target)
			val, failed := gr.getLastValue()
			if (failed == true){
				logMessage(" - no value was found")
			} else {
				logMessage(" - setting value %f for gauge: %+v", val, target)
				g.WithLabelValues(target).Set(val)
			}
		}
	}
	log.Printf("done collecting metrics\n\n")
}

func httpWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		collectMetrics()
		h.ServeHTTP(w, r)
	})
}

func main() {
	log.Println("Started Main")
	config = getConfig()

	// graphite.ssl.skiptls = config.SSLConfig.SkipTLS
	graphite.ssl = config.SSLConfig
	collectMetrics()

	http.Handle(getHTTPEndpoint(), httpWrapper(prometheus.Handler()))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html>
			<head><title>Graphite-Exporter</title></head>
			<body>
			<h1>Graphite Exporter</h1>
			<p><a href="` + HTTPEndpoint + `">Metrics</a></p>
			</body>
			</html>`))
	})
	log.Fatal(http.ListenAndServe(getHTTPPort(), nil))

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
