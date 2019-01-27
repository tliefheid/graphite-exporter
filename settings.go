package main

import (
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
}

// Config struct
type Config struct {
	GraphiteURL string   `yaml:"graphite"`
	Metrics     []Metric `yaml:"metrics"`
}

func getConfig() Config {
	log.Println("getting config")
	yml, err := ioutil.ReadFile("config.yml")
	check(err)
	c := Config{}
	config := Config{GraphiteURL: c.GraphiteURL}
	err = yaml.Unmarshal(yml, &c)
	check(err)
	for _, data := range c.Metrics {
		if data.GraphiteURL == "" {
			data.GraphiteURL = c.GraphiteURL
		}
		data.Name = trimAndReplace(data.Name)
		config.Metrics = append(config.Metrics, data)
	}
	return config
}
