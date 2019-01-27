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
		log.Println("collecting metrics")
		collectMetrics()
		h.ServeHTTP(w, r)
	})
}

func main() {
	log.Println("Started Main")
	config = getConfig()
	collectMetrics()

	http.Handle("/metrics", httpWrapper(prometheus.Handler()))
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
