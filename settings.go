package main

import (
	"fmt"
	"io/ioutil"
	"reflect"

	"gopkg.in/yaml.v2"
)

var (
	cfg       = config{}
	graphites = make(map[string]Graphite)
	targets   = make(map[string]Target)
)

// Ssl - ssl config
type Ssl struct {
	Credentials string `yaml:"credentials,omitempty"`
	Certificate string `yaml:"certificate,omitempty"`
	SkipTLS     bool   `yaml:"skip_tls,omitempty"`
}

// Graphite - graphite config
type Graphite struct {
	Name      string   `yaml:"name"`
	URL       string   `yaml:"url"`
	Labels    []string `yaml:"labels,omitempty"`
	Ssl       Ssl      `yaml:"ssl,omitempty"`
	Namespace string   `yaml:"namespace,omitempty"`
	Offset    int      `yaml:"offset,omitempty"`
}

// Target - target config
type Target struct {
	Name         string   `yaml:"name"`
	GraphiteName string   `yaml:"graphite"`
	Query        string   `yaml:"query"`
	Labels       []string `yaml:"labels,omitempty"`
	Namespace    string   `yaml:"namespace,omitempty"`
	Wildcards    []string `yaml:"wildcards,omitempty"`
}
type server struct {
	Port     int    `yaml:"port,omitempty"`
	Endpoint string `yaml:"endpoint,omitempty"`
	LogLevel string `yaml:"log_level,omitempty"`
}
type config struct {
	Graphite []Graphite `yaml:"graphite"`
	Server   server     `yaml:"server"`
	Targets  []Target   `yaml:"targets"`
}

// GraphiteResponse - response object for graphite queries
type GraphiteResponse struct {
	Target string `json:"target"`
	Tags   struct {
		Name string `json:"name"`
	} `json:"tags"`
	Datapoints [][]*float64 `json:"datapoints"`
}

func getConfig() {
	Log.Info("Getting config")
	yml, err := ioutil.ReadFile("./config/config.yml")
	check(err)

	err = yaml.Unmarshal(yml, &cfg)
	check(err)
	Log.Debugf("config\n%+v\n", cfg)
	for _, graphiteConfig := range cfg.Graphite {
		if graphiteConfig.Name == "" {
			panic("Config error: Graphite name can't be empty")
		}
		if graphiteConfig.URL == "" {
			panic("Config error: Graphite url can't be empty")
		}
		name := trimAndReplace(graphiteConfig.Name)
		_, found := graphites[name]
		if found {
			keys := reflect.ValueOf(graphites).MapKeys()
			panic(fmt.Sprintf("you already defined a graphite instance with name: %v\nalready defined names: %v", name, keys))
		}
		// trimAndReplace(graphiteConfig.Name)
		// no spaces or - in metric names
		// no spaces or - in namespaces
		trimAndReplaceRef(&graphiteConfig.Namespace)

		graphites[graphiteConfig.Name] = graphiteConfig
		graphiteConfig.init()

	}

	for _, target := range cfg.Targets {
		trimAndReplaceRef(&target.Name)
		trimAndReplaceRef(&target.Namespace)
		g, ok := graphites[target.GraphiteName]
		if !ok {
			s := fmt.Sprintf("searched for graphite (%v) but couldnt find it for target (%v)", target.GraphiteName, target.Name)
			panic(s)
		}
		target.init(g)
	}
}
func getHTTPEndpoint() string {
	endpoint := "/metrics"
	if cfg.Server.Endpoint != "" {
		endpoint = cfg.Server.Endpoint
	}
	return endpoint
}
func getHTTPPort() string {
	port := 8080
	if cfg.Server.Port != 0 && cfg.Server.Port > 0 {
		port = cfg.Server.Port
	}
	p := fmt.Sprintf("%v", port)
	return ":" + p
}
