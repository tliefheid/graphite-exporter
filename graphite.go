package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// type sslConfig struct {
// 	skiptls         bool
// 	credentials     string
// 	certificatePath string
// }

// Graphite struct
type Graphite struct {
	url string
	ssl SSLConfig
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
	client := getClient(g, url)
	gr := getResponse(client, url, g.ssl.Credentials)

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
	t := time.Now().Unix()
	// fmt.Printf("\ntime: %v\n", t)
	unix := time.Now().Unix() - 120
	// fmt.Printf("past: %v\n\n", unix)
	logMessage("   - now: %v, query time (now-offset): %v", t, unix)
	return fmt.Sprintf("%v", unix)
}

func getClient(graphite Graphite, url string) *http.Client {
	hasPrefix := strings.HasPrefix(url, "https://")
	if hasPrefix == true {
		// https client
		return createHTTPSClient(graphite.ssl)
	}
	// plain client
	logMessage("  - creating plain http client")
	return &http.Client{}

}

func createHTTPSClient(ssl SSLConfig) *http.Client {
	logMessage("  - creating HTTPS Client")
	var tls = &tls.Config{}
	if ssl.SkipTLS == true {
		// create client with insecureSkipVerify
		logMessage("    - setting InsecureSkipVerify: true")
		tls.InsecureSkipVerify = ssl.SkipTLS
	}

	if ssl.CertificatePath != "" {
		// read certificate
		logMessage("     - setting ssl certificate")
		exists, _ := exists(ssl.CertificatePath)
		if exists == true {
			dat, err := ioutil.ReadFile(ssl.CertificatePath)
			check(err)
			rootPEM := string(dat)

			certPool := x509.NewCertPool()
			ok := certPool.AppendCertsFromPEM([]byte(rootPEM))
			if !ok {
				panic("failed to parse root certificate")
			}
			tls.RootCAs = certPool
		} else {
			logMessage("     - certificate not found")
		}

	}

	transport := &http.Transport{TLSClientConfig: tls}
	return &http.Client{Transport: transport}

}

func getResponse(client *http.Client, url string, credentials string) []GraphiteResponse {

	req, err := http.NewRequest("GET", url, nil)
	check(err)

	if credentials != "" {
		// set header
		encoded := base64.StdEncoding.EncodeToString([]byte(credentials))
		req.Header.Add("Authorization", "Basic "+encoded)
	}

	response, err := client.Do(req)
	check(err)
	// parsing response
	bytes := getBytes(response)
	logMessage("   - raw response: %+v\n", string(bytes))

	var gr []GraphiteResponse
	json.Unmarshal(bytes, &gr)
	// check(marshallError)
	logMessage("   - graphite response: %+v", gr)
	return gr
}

func (gr GraphiteResponse) getLastValue() (float64, bool) {
	returnValue := float64(-1)
	for i := len(gr.Datapoints) - 1; i >= 0; i-- {
		dp := gr.Datapoints[i]
		value := dp[0]
		time := dp[1]
		logMessage("   - %v | value: %v, time: %v", i, value, time)
		if value != 0 {
			returnValue = value
			return value, false
		}
	}

	return returnValue, true
}
