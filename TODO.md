# ToDo

- ~~make http port customizable~~
- custom seperator for labels
- add labels at global level
- ~~make namespace customizable, both globally and per metric~~
- wildcard labels
- add tests
- ~~add Dockerfile~~
- ~~make metric endpoint customizable~~

## Wildcard label idea

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