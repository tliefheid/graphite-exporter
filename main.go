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
)

func collectMetrics() {
	log.Println("collecting metrics")
	for _, m := range config.Metrics {
		g := getGauge(m)
		respSlice := graphite.getMetric(m)
		for _, gr := range respSlice {
			target := trimAndReplace(gr.Target)
			val := gr.getLastValue()
			g.WithLabelValues(target).Set(val)
		}
	}
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
	collectMetrics()

	http.Handle(getHTTPEndoint(), httpWrapper(prometheus.Handler()))
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
