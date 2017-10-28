package createservice_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/pivotal-cloudops/cf-sli/cf_wrapper/cf_wrapperfakes"
	"github.com/pivotal-cloudops/cf-sli/config"
	"github.com/pivotal-cloudops/cf-sli/sli_executor"
	"github.com/pivotal-cloudops/cf-sli/ssl_sli/create_service"
	datadoghelpers "github.com/pivotal-cloudops/cf-sli/ssl_sli/datadog_helpers"
)

var _ = Describe("CreateService", func() {

	var (
		fakeCf      *cf_wrapperfakes.FakeCfWrapperInterface
		sliExecutor sli_executor.SliExecutor
		config      config.Config
		datadogInfo datadoghelpers.DatadogInfo
	)

	BeforeEach(func() {
		fakeCf = new(cf_wrapperfakes.FakeCfWrapperInterface)
		sliExecutor = sli_executor.NewSliExecutor(fakeCf)
		config.LoadConfig("../../fixtures/config_test.json")
		datadogInfo = datadoghelpers.DatadogInfo{
			DatadogAPIKey:  "fakeKey",
			DatadogAppKey:  "fakeAppKey",
			DeploymentName: "fakeDep",
			Metric:         "some-metric",
		}
	})

	Context("#SLI", func() {
		It("Posts to datadog with status 1", func() {
			err := createservice.SLI(config, sliExecutor, datadogInfo)
			output := createservice.CaptureStdout(func() { createservice.SLI(config, sliExecutor, datadogInfo) })

			Expect(err).To(Equal(""))
			Expect(output).To(ContainSubstring("Create status: 1 for metric "))
		})

		Context("When preparing the SLI fails", func() {
			It("Returns the error", func() {
				fakeCf.StubFailingCF("target")
				err := createservice.SLI(config, sliExecutor, datadogInfo)

				Expect(err).To(ContainSubstring("Running CF command failed: target"))
			})
		})

		Context("When creating the service fails", func() {
			It("Posts to datadog with status 0", func() {
				fakeCf.StubFailingCF("create-service")
				err := createservice.SLI(config, sliExecutor, datadogInfo)
				output := createservice.CaptureStdout(func() { createservice.SLI(config, sliExecutor, datadogInfo) })

				Expect(err).To(ContainSubstring("Running CF command failed: create-service"))
				Expect(output).To(ContainSubstring("Create status: 0 for metric "))
			})
		})

		Context("When the cleanup fails", func() {
			It("Returns the error, but does not affect posting to datadog", func() {
				fakeCf.StubFailingCF("delete-service")
				err := createservice.SLI(config, sliExecutor, datadogInfo)
				output := createservice.CaptureStdout(func() { createservice.SLI(config, sliExecutor, datadogInfo) })

				Expect(err).To(ContainSubstring("Running CF command failed: delete-service"))
				Expect(output).To(ContainSubstring("Create status: 1 for metric "))
			})
		})
	})
})
