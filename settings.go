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

// Config struct
type Config struct {
	GraphiteURL  string   `yaml:"graphite"`
	HTTPPort     int      `yaml:"http_port"`
	HTTPEndpoint string   `yaml:"http_endpoint"`
	Namespace    string   `yaml:"namespace"`
	Metrics      []Metric `yaml:"metrics"`
}

func getConfig() Config {
	log.Println("getting config")
	yml, err := ioutil.ReadFile("config.yml")
	check(err)
	c := Config{}
	config := Config{GraphiteURL: c.GraphiteURL}
	err = yaml.Unmarshal(yml, &c)
	check(err)

	if c.HTTPEndpoint != "" {
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
	return config
}
