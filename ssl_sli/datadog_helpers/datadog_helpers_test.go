package datadoghelpers_test

import (
	"io/ioutil"
	"log"

	datadoghelpers "github.com/pivotal-cloudops/cf-sli/ssl_sli/datadog_helpers"

	httpmock "github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DatadogHelpers", func() {
	BeforeEach(func() {
		log.SetOutput(ioutil.Discard)
	})

	It("makes a POST request to Datadog", func() {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("POST", "https://app.datadoghq.com/api/v1/series?api_key=blah&application_key=thing",
			httpmock.NewStringResponder(200, "OK"))

		datadogInfo := datadoghelpers.DatadogInfo{
			DatadogAPIKey:  "blah",
			DatadogAppKey:  "thing",
			DeploymentName: "bosher",
			Metric:         "ssl-sli",
		}
		status := datadoghelpers.PostToDatadog(1, datadogInfo)
		Expect(status).To(Equal("200"))
	})
})
