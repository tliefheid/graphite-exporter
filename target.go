package main

import (
	"strconv"
	"strings"

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
	// TODO: check if labels give errors and no duplicate labels
	constantLabels := append(t.Labels, graphite.Labels...)

	ns := t.getNamespaces()

	wildcardValues := getValuesFromArray(t.Wildcards, ":")
	customLabels := append(wildcardValues, "target")

	gauge = buildPrometheusGauge(t.Name, ns, constantLabels, customLabels)
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
		customLabels := t.getWildcardValues(data)
		Log.Debugf("Custom labels: %v\n", customLabels)

		val, fail := data.getLastValue()
		Log.Infof("getLastValue: failed?: '%v', value: %v for target: '%v'\n", fail, val, data.Target)
		if fail {
			Log.Info("no data found")
		} else {
			Log.Debug("setting value to gauge")
		}
		gauge.WithLabelValues(customLabels...).Set(val)
	}
}

func (t Target) getWildcardValues(data GraphiteResponse) []string {
	var output []string
	target := trimAndReplace(data.Target)
	sliced := strings.Split(target, ".")
	keys := getKeysFromArray(t.Wildcards, ":")
	for _, k := range keys {
		i, _ := strconv.Atoi(k)
		val := sliced[i]
		trimAndReplaceRef(&val)
		output = append(output, val)
	}
	output = append(output, target)
	return output
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
