package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var (
	client = http.Client{}
)

func (g Graphite) init() {
	g.createClient()
}

func (g Graphite) createClient() {
	Log.Info("creating HTTPS Client")
	var tls = tls.Config{}
	g.setTLSSkip(&tls)
	g.setClientCertificate(&tls)
	transport := &http.Transport{TLSClientConfig: &tls}
	client = http.Client{Transport: transport}
}

func (g Graphite) setTLSSkip(tls *tls.Config) {
	if g.Ssl.SkipTLS == true {
		Log.Info("setting InsecureSkipVerify: true")
		tls.InsecureSkipVerify = g.Ssl.SkipTLS
	}
}

func (g Graphite) setClientCertificate(tls *tls.Config) {
	path := g.Ssl.Certificate
	if path != "" {
		Log.Info("setting ssl certificate")
		exists, _ := exists(path)
		if exists == true {
			dat, err := ioutil.ReadFile(path)
			check(err)
			rootPEM := string(dat)

			certPool := x509.NewCertPool()
			ok := certPool.AppendCertsFromPEM([]byte(rootPEM))
			if !ok {
				panic("failed to parse root certificate")
			}
			tls.RootCAs = certPool
		} else {
			Log.Warning("certificate not found")
		}

	}
}
func (g Graphite) query(q string) []GraphiteResponse {
	// build url
	// build req
	// check for extra headers
	// do request
	// return response
	url := g.buildURL(q)
	req, err := http.NewRequest("GET", url, nil)
	check(err)
	// make more generic, create option to set multiple headers
	setHeader(req, "Authorization", g.Ssl.Credentials)

	res, err := client.Do(req)
	check(err)

	// parsing response
	bytes := getBytes(res)
	Log.Debugf("raw response: %+v\n", string(bytes))

	var gr []GraphiteResponse
	json.Unmarshal(bytes, &gr)

	Log.Debugf("graphite response: %+v", gr)
	return gr
}

func (g Graphite) buildURL(query string) string {
	url := g.URL
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}
	url += "render?target=" + query + "&format=json&from=" + getFromTime(g.Offset)
	Log.Infof("query url: %v", url)
	return url
}

func setHeader(req *http.Request, key string, value string) {
	encoded := base64.StdEncoding.EncodeToString([]byte(value))
	req.Header.Add(key, "Basic "+encoded)
}

func getFromTime(offset int) string {
	if offset <= 0 {
		offset = 60
	}
	t := time.Now().Unix()
	unix := time.Now().Unix() - int64(offset)
	Log.Infof("now: %v, query time (now-offset(%v)): %v", t, offset, unix)
	return fmt.Sprintf("%v", unix)
}

func (gr GraphiteResponse) getLastValue() (float64, bool) {
	for i := len(gr.Datapoints) - 1; i >= 0; i-- {
		dp := gr.Datapoints[i]
		if dp[0] != nil {
        value := *dp[0]
        time := *dp[1]
        Log.Infof("value: %v, time: %v", value, time)
			  return value, false
		}
	}

	return float64(-1), true
}
