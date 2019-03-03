package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

func findGauge(name string) (bool, prometheus.GaugeVec) {
	for k, v := range gauges {
		if k == name {
			logMessage(" - found a gauge")
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
			Namespace:   m.Namespace,
			ConstLabels: lblMap,
			Name:        m.Name,
		},
		[]string{
			"target",
		},
	)
	prometheus.MustRegister(g)
	gauges[m.Name] = *g
	logMessage("   - new gauge build complete: %s", m.Name)
	return *g
}

func getGauge(m Metric) prometheus.GaugeVec {
	logMessage("   - getGauge() for: %s", m.Name)
	name := trimAndReplace(m.Name)
	found, g := findGauge(name)
	if found {
		return g
	}
	return buildGauge(m)

}
