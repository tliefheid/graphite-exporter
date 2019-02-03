# Graphite Exporter

This is a exporter for prometheus. Unlike the official exporter, it queries graphite instead of receiving graphite metrics.

You provide graphite queries to the exporter. If you call the metrics endpoint, it queries graphite and exposes the results.

## Run in Docker

```Shell
docker run -d \
-v /path/to/config.yml:/app/config.yml:ro \
-p 8080:8080 \
tomldev/graphite-exporter
```

## Configuration

You can define a global graphite instance, which you can override for each metric

minimal config:

```YAML
---
graphite: http://graphite.instance.com:1234/
metrics:
  - name: foo
    query: some.graphite.query.*
```

extended config:

```YAML
---
graphite: http://graphite.instance.com/
http_port: 9009 # default: 8080
http_endpoint: /custom/metric/endpoint # default: /metrics
namespace: custom_namespace # default: graphite_exporter

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
