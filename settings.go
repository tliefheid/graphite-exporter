package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// Metric single metric config
type Metric struct {
	Name        string   `yaml:"name"`
	Query       string   `yaml:"query"`
	Labels      []string `yaml:"labels,omitempty"`
	GraphiteURL string   `yaml:"graphite,omitempty"`
	Namespace   string   `yaml:"namespace,omitempty"`
}

// SSLConfig ssl configuration
type SSLConfig struct {
	SkipTLS         bool   `yaml:"skip_tls,omitempty"`
	CertificatePath string `yaml:"certificate_path,omitempty"`
	Credentials     string `yaml:"credentials,omitempty"`
}

// Config struct
type Config struct {
	GraphiteURL  string    `yaml:"graphite"`
	HTTPPort     int       `yaml:"http_port"`
	HTTPEndpoint string    `yaml:"http_endpoint"`
	Namespace    string    `yaml:"namespace"`
	SSLConfig    SSLConfig `yaml:"ssl"`
	Debug        bool      `yaml:"debug"`
	Metrics      []Metric  `yaml:"metrics"`
}

// SkipTLS      bool     `yaml:"skip_tls"`

func getConfig() Config {
	log.Println("getting config")
	yml, err := ioutil.ReadFile("config.yml")
	check(err)
	// marshalled config from file
	c := Config{}

	// final config
	config := Config{GraphiteURL: c.GraphiteURL}
	
	err = yaml.Unmarshal(yml, &c)
	check(err)

	if c.SSLConfig.SkipTLS == true {
		config.SSLConfig.SkipTLS = true
	} else {
		config.SSLConfig.SkipTLS = false
	}
	config.SSLConfig = c.SSLConfig
	if c.Debug == true {
		DebugLogging = true
	}

	if c.HTTPEndpoint != "" {
		// overwriting global setting
		HTTPEndpoint = c.HTTPEndpoint
		log.Println("  - Setting metrics endpoint to: " + HTTPEndpoint)
	}

	if c.HTTPPort != 0 && c.HTTPPort > 0 {
		HTTPPort = c.HTTPPort
		log.Println("  - Setting server port to: " + fmt.Sprintf("%v", HTTPPort))
	}

	if c.Namespace != "" {
		namespace = trimAndReplace(c.Namespace)
		log.Println("  - Setting custom global namespace to: " + namespace)
	}

	for _, data := range c.Metrics {
		if data.GraphiteURL == "" {
			data.GraphiteURL = c.GraphiteURL
		}
		if data.Namespace == "" {
			data.Namespace = namespace
		}
		data.Namespace = trimAndReplace(data.Namespace)

		data.Name = trimAndReplace(data.Name)
		config.Metrics = append(config.Metrics, data)
	}
	logMessage("config\n%+v\n", config)
	return config
}
