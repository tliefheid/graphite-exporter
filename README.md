# Graphite Exporter

This is a exporter for prometheus. Unlike the official exporter, it queries graphite instead of receiving graphite metrics.

You provide graphite queries to the exporter. If you call the metrics endpoint, it queries graphite and exposes the results.

## config.yml

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
    namespace: metric_specific_namespace
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

## ToDo

- make http port customizable
- custom seperator for labels
- add labels at global level
- make namespace customizable, both globally and per metric
- wildcard labels
- add tests
- add Dockerfile
- make metric endpoint customizable

### Wildcard label

Lets say a sensor has 2 values in graphite, humidity and temperature.

```txt
sensor.values.humidity = 36.0
sensor.values.temperature = 18.0
```

If you do a query like `sensor.values.*` you might want the result to be:

```Go
graphite_exporter_sensor{target="sensor.values.humidity", type="humidity"} 36.0
graphite_exporter_sensor{target="sensor.values.temperature", type="temperature"} 18.0
```