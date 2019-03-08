# Graphite Exporter

This is a exporter for prometheus. Unlike the official exporter, it queries graphite instead of receiving graphite metrics.

You provide graphite queries to the exporter. If you call the metrics endpoint, it queries graphite and exposes the results.

## Run in Docker

```Shell
docker run -d \
-v /path/to/config.yml:/app/config.yml:ro \
-v /path/to/certificate/my-cert:/etc/certs/root.cer \
-p 8080:8080 \
tomldev/graphite-exporter
```

or use a compose file:

```YAML
version: '3.3'

networks:
  networkname:
    external: true

services:
  graphiteexporter:
    image: tomldev/graphite-exporter:v1.3.0
    networks:
      - networkname
    ports:
      - "9999:8080"
    volumes:
      - ./config.yml:/app/config.yml
      - ./certs/my-cert.cer:/etc/certs/root.cer
```

Use docker-compose (`docker-compose up -d`) or a stack deploy to a swarm cluster (`docker stack deploy --compose-file docker-compose.yml STACKNAME`)

## Configuration

**minimal config:**

```YAML
---
graphite: http://graphite.instance.com:1234/
metrics:
  - name: foo
    query: some.graphite.query.*
```

**extended config:**

```YAML
---
graphite: http://graphite.instance.com/
http_port: 9009 # default: 8080
http_endpoint: /custom/metric/endpoint # default: /metrics
namespace: custom_namespace # default: graphite_exporter
skip_tls: true # deprecated since 1.3.0
debug: true # default false

ssl:
  credentials: 'username:password' # this will generate authorization header with 'Basic <base64 encoded credentials>'
  certificate_path: '/etc/certs/root.cer' # path to a certificate
  skip_tls: false # skip tls validation

metrics:
  - name: foo
    query: some.graphite.query.*
    labels:
      - 'label1: value1'
  - name: bar
    query: some.other.graphite.query
    labels:
      - 'label1: value1'
      - 'label2: value2'
  - name: external graphite
    namespace: metric_specific_namespace # overwrite at metric level
    graphite: http://external.graphite.instance.com/
    query: external.graphite.query
```

**explanation:**

- graphite: Global graphite connection.
- http_port: The port on which the metrics will be exposed.
- http_endpoint: On which endpoint you want to expose your metrics
- namespace: global metric name prefix
- skip_tls: skip tls verification on your graphite instance (deprecated since 1.3.0)
- debug: debug logging
- ssl:
  - credentials: when provided, the request to graphite will be send with an Authorization header with `Basic: <token>`. The token will be an base64 encoded string of the credentials
  - certificate_path: when provided, the request to graphite will be send with the specified certificate
  - skip_tls: skip tls verification on your graphite instance
- metrics:
  - name: name of the metric
  - query: the graphite query
  - labels: add custom key:value labels to your metric
  - graphite: overwrite graphite connection for a metric
  - namespace: overwrite global namespace for a metric

All spaces in the names and labels will be trimmed and the remaining spaces will be replaced by an `_`

For the labels you need to use an `:` as seperator

You can also use the graphite query wildcard. The query is added to the target label.

## Result

```Go
graphite_exporter_foo{label1="value1", target="some.graphite.query.query1"} 10.0
graphite_exporter_foo{label1="value1", target="some.graphite.query.query2"} 20.0
graphite_exporter_bar{label1="value1", label2="value2", target="some.other.graphite.query"} 42.0
graphite_exporter_external_graphite{target="external.graphite.query"} 65.0
```
