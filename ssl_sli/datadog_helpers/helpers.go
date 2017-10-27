package datadoghelpers

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type DatadogInfo struct {
	DatadogAPIKey  string
	DatadogAppKey  string
	DeploymentName string
	Metric         string
}

func PostToDatadog(result int, datadogInfo DatadogInfo) {
	datadogURL := "https://app.datadoghq.com/api/v1/series?api_key=" + datadogInfo.DatadogAPIKey + "&application_key=" + datadogInfo.DatadogAppKey
	currentTime := time.Now()

	body := createPOSTBody(result, datadogInfo, currentTime)
	req, _ := http.NewRequest("POST", datadogURL, body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	fmt.Println("Time: ", currentTime, "Posting to datadog: ", datadogInfo.Metric, result)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("error posting to Datadog: ", err)
		os.Exit(1)
	}

	defer resp.Body.Close()
	fmt.Println("Posted to datadog: ", resp.Status)
}

func createPOSTBody(result int, datadogInfo DatadogInfo, currentTime time.Time) *strings.Reader {
	epochTime := currentTime.Unix()
	body := fmt.Sprintf(`
		{ "series" :
         [
					 {"metric":"%s",
						"points":[[%v, %d]],
						"type":"gauge",
						"tags":["sli","deployment:%s"]}
        ]
		}`, datadogInfo.Metric, epochTime, result, datadogInfo.DeploymentName)

	return strings.NewReader(body)
}
