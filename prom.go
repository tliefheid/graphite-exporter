package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

func buildPrometheusGauge(name string, ns string, constantLabels []string, customLabels []string) prometheus.GaugeVec {
	// create map from the labels
	lblMap := make(map[string]string)
	for _, lbl := range constantLabels {
		k, v := getKeyValue(lbl, ":")
		lblMap[k] = v
	}

	g := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace:   ns,
			ConstLabels: lblMap,
			Name:        name,
		},
		customLabels,
	)

	Log.Debug("created new gauge with name: %v", name)
	prometheus.MustRegister(g)
	return *g
}
