# cf-sli

Service Level Indicator monitors for [Pivotal Web Services](http://run.pivotal.io/)

Calculates percentage of pushes successful over a period of time

## Prerequisites
* Stable [CF CLI](https://github.com/cloudfoundry/cli/releases)
* Config file (see sample below)
* Datadog Account (if you would like to send metrics somewhere)

## Sample config file
```
cat > .config
{
  "api": "api.example.com",
  "user": "cf-sli",
  "pass": "abc123",
  "org": "cf-sli",
  "space": "dev"
}
```

## Output format
```json
{
  "app_start_time_in_sec": "Number", 
  "app_stop_time_in_sec": "Number",
  "app_start_status": "Number", 
  "app_stop_status": "Number"
}
```

_inspired by Google SRE's [Service Level Objectives](https://landing.google.com/sre/book/chapters/service-level-objectives.html)_
