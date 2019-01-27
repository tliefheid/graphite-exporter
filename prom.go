package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

func findGauge(name string) (bool, prometheus.GaugeVec) {
	for k, v := range gauges {
		if k == name {
			return true, v
		}
	}
	return false, prometheus.GaugeVec{}
}

func buildGauge(m Metric) prometheus.GaugeVec {
	labels := m.Labels
	lblMap := make(map[string]string)
	for _, lbl := range labels {
		key, val := getKeyValue(lbl, ":")
		lblMap[key] = val
	}
	g := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace:   "graphite_exporter",
			ConstLabels: lblMap,
			Name:        m.Name,
		},
		[]string{
			"target",
		},
	)
	prometheus.MustRegister(g)
	gauges[m.Name] = *g
	return *g
}

func getGauge(m Metric) prometheus.GaugeVec {
	// fmt.Printf("gauges: %v\n", gauges)
	name := trimAndReplace(m.Name)
	found, g := findGauge(name)
	if found {
		return g
	}
	return buildGauge(m)

}
