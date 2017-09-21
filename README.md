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
  "space": "dev",
  "domain": "example.com"
  "domain": "example.com",
  "stack": "cflinuxfs2"
}
```

## Run, extract metrics, and post to datadog
```bash
output=$(cf-sli)
stats=$(echo -e "$output" | tail -n 1) # extract last line to ignore cf cli output
app_route=$(echo $stats | jq .app_route)
app_start_time=$(echo $stats | jq .app_start_time_in_sec)
app_stop_time=$(echo $stats | jq .app_stop_time_in_sec)
app_start_status=$(echo $stats | jq .app_start_status)
app_stop_status=$(echo $stats | jq .app_stop_status)

currenttime=$(date +%s)

echo "Sending metrics to datadog..."

post_to_datadog () {
  local metric=$1
  local data_point_time=$2
  local data_point_value=$3
  local tag=$4

  curl -X POST -H "Content-type: application/json" \
    -d "{ \"series\" :
           [{\"metric\": \"$metric\",
            \"points\": [[$data_point_time, $data_point_value]],
            \"type\": \"gauge\",
            \"tags\": [\"$tag\"]
          }]
        }" \
  "https://app.datadoghq.com/api/v1/series?api_key=${DATADOG_API_KEY}"
}

post_to_datadog "cloudops_tools.cf_sli.app_start_time_in_sec" $currenttime $app_start_time "deployment:${BOSH_DEPLOYMENT_NAME}"
post_to_datadog "cloudops_tools.cf_sli.app_stop_time_in_sec" $currenttime $app_stop_time "deployment:${BOSH_DEPLOYMENT_NAME}"
post_to_datadog "cloudops_tools.cf_sli.app_start_status" $currenttime $app_start_status "deployment:${BOSH_DEPLOYMENT_NAME}"
post_to_datadog "cloudops_tools.cf_sli.app_stop_status" $currenttime $app_stop_status "deployment:${BOSH_DEPLOYMENT_NAME}"
```

_inspired by Google SRE's [Service Level Objectives](https://landing.google.com/sre/book/chapters/service-level-objectives.html)_
