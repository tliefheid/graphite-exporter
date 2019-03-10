package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	graphite = Graphite{}
	gauge    = prometheus.GaugeVec{}
)

func (t Target) init(g Graphite) {
	// create ref to graphite
	// build gauge
	graphite = g
	// TODO: check if labels give errors
	labels := append(t.Labels, graphite.Labels...)
	ns := t.getNamespaces()
	gauge = buildPrometheusGauge(t.Name, ns, labels)
}

func (t Target) getMetrics() {
	Log.Noticef("collecting metrics for: %v", t.Name)
	// query graphite
	// generate response
	// set value of gauge
	res := graphite.query(t.Query)

	if len(res) < 1 {
		Log.Noticef("no data retreived from graphite for Target: %v", t.Name)
	}

	for _, data := range res {
		target := trimAndReplace(data.Target)
		val, fail := data.getLastValue()
		Log.Infof("value: %v or failed (%v) for target: %v\n", val, fail, target)
		if fail {
			Log.Info("no data found")
		} else {
			Log.Debug("setting value to gauge")
			gauge.WithLabelValues(target).Set(val)
		}
	}
	// fmt.Printf("response: %v\n", res)
}

func (t Target) getNamespaces() string {
	defaultNs := "graphite_exporter"
	ns := defaultNs
	if graphite.Namespace != "" {
		ns = graphite.Namespace
	}
	if t.Namespace != "" {
		if ns == defaultNs {
			ns = t.Namespace
		} else {
			ns += "_" + t.Namespace
		}
	}
	return ns
}
