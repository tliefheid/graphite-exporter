package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// Graphite strucs
type Graphite struct {
	url     string
	skiptls bool
}

// GraphiteResponse struct for graphite json response
type GraphiteResponse struct {
	Target string `json:"target"`
	Tags   struct {
		Name string `json:"name"`
	} `json:"tags"`
	Datapoints [][]float64 `json:"datapoints"`
}

func (g Graphite) getMetric(m Metric) []GraphiteResponse {
	logMessage("   - getting metrics for: %s", m.Name)
	url := buildURL(m)

	gr := getResponse(url, g.skiptls)

	if len(gr) == 0 {
		log.Println("   - no data found in graphite for: " + m.Name)
	}
	return gr
}

func buildURL(m Metric) string {
	url := ""
	url += m.GraphiteURL
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}
	url += "render?target="
	url += m.Query
	url += "&format=json&from="
	url += getFromTime()
	logMessage("   - build graphite url: %s", url)
	return url
}

func getFromTime() string {
	unix := time.Now().Unix() - 30
	return fmt.Sprintf("%v", unix)
}

func getResponse(url string, skip bool) []GraphiteResponse {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skip},
	}

	client := &http.Client{Transport: tr}
	response, err := client.Get(url)

	// response, err := http.Get(url)
	check(err)

	bytes := getBytes(response)
	var gr []GraphiteResponse
	json.Unmarshal(bytes, &gr)
	logMessage("   - graphite response: %+v", gr)
	return gr
}

func (gr GraphiteResponse) getLastValue() float64 {
	dp := gr.Datapoints
	l := len(dp)

	last := dp[l-1]
	val := last[0]
	return val
}
