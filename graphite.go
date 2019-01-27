package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Graphite strucs
type Graphite struct {
	url string
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
	url := buildURL(m)
	gr := getResponse(url)
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
	return url
}

func getFromTime() string {
	unix := time.Now().Unix() - 30
	return fmt.Sprintf("%v", unix)
}

func getResponse(url string) []GraphiteResponse {
	response, err := http.Get(url)
	check(err)

	bytes := getBytes(response)
	var gr []GraphiteResponse
	json.Unmarshal(bytes, &gr)
	return gr
}

func (gr GraphiteResponse) getLastValue() float64 {
	dp := gr.Datapoints
	l := len(dp)

	last := dp[l-1]
	val := last[0]
	return val
}
