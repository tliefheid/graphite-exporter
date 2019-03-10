package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

func buildPrometheusGauge(name string, ns string, labels []string) prometheus.GaugeVec {
	// create map from the labels
	lblMap := make(map[string]string)
	for _, lbl := range labels {
		k, v := getKeyValue(lbl, ":")
		lblMap[k] = v
	}

	g := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace:   ns,
			ConstLabels: lblMap,
			Name:        name,
		},
		[]string{
			"target",
		},
	)

	Log.Debug("created new gauge with name: %v", name)
	prometheus.MustRegister(g)
	return *g
}
