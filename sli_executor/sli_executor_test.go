package sli_executor_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cloudops/cf-sli/cf_wrapper/cf_wrapperfakes"
	"github.com/pivotal-cloudops/cf-sli/config"
	"github.com/pivotal-cloudops/cf-sli/sli_executor"
	. "github.com/tjarratt/gcounterfeiter"
)

var _ = Describe("SliExecutor", func() {
	var expected_timeout = "60"
	var expected_api_calls = []string{"api", "fake_api"}
	var expected_auth_calls = []string{"auth", "fake_user", "fake_pass"}
	var expected_target_calls = []string{"target", "-o", "fake_org", "-s", "fake_space"}
	var expected_push_calls = []string{"push", "-p", "./fake_path", "-b", "fake_buildpack", "fake_app_name", "-d", "fake_domain", "--no-start", "-t", expected_timeout}
	var expected_start_calls = []string{"start", "fake_app_name"}
	var expected_stop_calls = []string{"stop", "fake_app_name"}
	var expected_delete_calls = []string{"delete", "fake_app_name", "-f", "-r"}
	var expected_logout_calls = []string{"logout"}
	var expected_logs_calls = []string{"logs", "fake_app_name", "--recent"}

	var (
		fakeCf *cf_wrapperfakes.FakeCfWrapperInterface
		sli    sli_executor.SliExecutor
		config config.Config
	)
	BeforeEach(func() {
		fakeCf = new(cf_wrapperfakes.FakeCfWrapperInterface)
		sli = sli_executor.NewSliExecutor(fakeCf)
		config.LoadConfig("../fixtures/config_test.json")
	})

	Context("#Prepare", func() {
		It("returns nil if cf command executes successfully", func() {
			err := sli.Prepare("fake_api", "fake_user", "fake_pass", "fake_org", "fake_space")
			Expect(err).NotTo(HaveOccurred())

			Expect(fakeCf).To(HaveReceived("RunCF").With(expected_api_calls))
			Expect(fakeCf).To(HaveReceived("RunCF").With(expected_auth_calls))
			Expect(fakeCf).To(HaveReceived("RunCF").With(expected_target_calls))
		})

		It("returns err when cf api fails", func() {
			fakeCf.StubFailingCF("api")
			err := sli.Prepare("fake_api", "fake_user", "fake_pass", "fake_org", "fake_space")
			Expect(err).To(HaveOccurred())
		})

		It("returns err when cf auth fails", func() {
			fakeCf.StubFailingCF("auth")
			err := sli.Prepare("fake_api", "fake_user", "fake_pass", "fake_org", "fake_space")
			Expect(err).To(HaveOccurred())
		})

		It("returns err when cf target fails", func() {
			fakeCf.StubFailingCF("target")
			err := sli.Prepare("fake_api", "fake_user", "fake_pass", "fake_org", "fake_space")
			Expect(err).To(HaveOccurred())
		})
	})

	Context("#PushAndStartSli", func() {
		It("Push the Sli app with --no-start and starts it", func() {
			elapsed_time, err := sli.PushAndStartSli("fake_app_name", "fake_buildpack", "fake_domain", "./fake_path", "60")
			Expect(err).NotTo(HaveOccurred())
			Expect(elapsed_time).ToNot(Equal(0))

			Expect(fakeCf).To(HaveReceived("RunCF").With(expected_push_calls))
			Expect(fakeCf).To(HaveReceived("RunCF").With(expected_start_calls))
		})

		It("Returns error when cf push fails", func() {
			fakeCf.StubFailingCF("push")
			elapsed_time, err := sli.PushAndStartSli("fake_app_name", "fake_buildpack", "fake_domain", "./fake_path", "60")
			Expect(err).To(HaveOccurred())
			Expect(elapsed_time).To(Equal(time.Duration(0)))
		})

		It("Returns error when cf start fails", func() {
			fakeCf.StubFailingCF("start")
			elapsed_time, err := sli.PushAndStartSli("fake_app_name", "fake_buildpack", "fake_domain", "./fake_path", "60")
			Expect(err).To(HaveOccurred())
			Expect(elapsed_time).To(Equal(time.Duration(0)))
		})
	})

	Context("#StopSli", func() {
		It("Start the Sli app", func() {
			elapsed_time, err := sli.StopSli("fake_app_name")
			Expect(err).NotTo(HaveOccurred())

			Expect(elapsed_time).ToNot(Equal(time.Duration(0)))

			Expect(fakeCf).To(HaveReceived("RunCF").With(expected_stop_calls))
		})

		It("Returns error when cf stop fails", func() {
			fakeCf.StubFailingCF("stop")
			elapsed_time, err := sli.StopSli("fake_app_name")
			Expect(err).To(HaveOccurred())
			Expect(elapsed_time).To(Equal(time.Duration(0)))
		})
	})

	Context("#CleanupSli", func() {
		It("delete the Sli app and logs out", func() {
			err := sli.CleanupSli("fake_app_name")
			Expect(err).NotTo(HaveOccurred())

			Expect(fakeCf).To(HaveReceived("RunCF").With(expected_delete_calls))
			Expect(fakeCf).To(HaveReceived("RunCF").With(expected_logout_calls))
		})

		It("Returns error when cf delete fails, and it logs out", func() {
			fakeCf.StubFailingCF("delete")
			err := sli.CleanupSli("fake_app_name")
			Expect(err).To(HaveOccurred())

			Expect(fakeCf).To(HaveReceived("RunCF").With(expected_logout_calls))
		})

		It("Returns error when cf logout fails", func() {
			fakeCf.StubFailingCF("logout")
			err := sli.CleanupSli("fake_app_name")
			Expect(err).To(HaveOccurred())
		})
	})

	Context("#RunTest", func() {
		It("Login, push the app, returns the start and stop times and status, and cleanup", func() {
			result, err := sli.RunTest("fake_app_name", "fake_buildpack", "./fake_path", config)
			Expect(err).NotTo(HaveOccurred())

			// Login and target to the org and space
			Expect(fakeCf).To(HaveReceived("RunCF").With(expected_api_calls))
			Expect(fakeCf).To(HaveReceived("RunCF").With(expected_auth_calls))
			Expect(fakeCf).To(HaveReceived("RunCF").With(expected_target_calls))

			// Push, start, and stop the app
			Expect(fakeCf).To(HaveReceived("RunCF").With(expected_push_calls))
			Expect(fakeCf).To(HaveReceived("RunCF").With(expected_start_calls))
			Expect(fakeCf).To(HaveReceived("RunCF").With(expected_stop_calls))

			Expect(result.StartTime).ToNot(Equal(0))
			Expect(result.StopTime).ToNot(Equal(0))
			Expect(result.StartStatus).To(Equal(1))
			Expect(result.StopStatus).To(Equal(1))

			// Cleanup and logout
			Expect(fakeCf).To(HaveReceived("RunCF").With(expected_delete_calls))
			Expect(fakeCf).To(HaveReceived("RunCF").With(expected_logout_calls))
		})

		Context("When something in the prepare step fails", func() {
			It("Cleans up the app", func() {
				fakeCf.StubFailingCF("api")
				sli.RunTest("fake_app_name", "fake_buildpack", "./fake_path", config)

				expected_delete_calls := []string{"delete", "fake_app_name", "-f", "-r"}
				Expect(fakeCf).To(HaveReceived("RunCF").With(expected_delete_calls))

				Expect(fakeCf).To(HaveReceived("RunCF").With(expected_logout_calls))
			})

			It("Returns an error from CF", func() {
				fakeCf.StubFailingCF("auth")
				_, err := sli.RunTest("fake_app_name", "fake_buildpack", "./fake_path", config)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Running CF command failed:"))
			})

			It("Does not record time and status", func() {
				fakeCf.StubFailingCF("target")
				result, _ := sli.RunTest("fake_app_name", "fake_buildpack", "./fake_path", config)

				Expect(result.StartTime).To(Equal(time.Duration(0)))
				Expect(result.StopTime).To(Equal(time.Duration(0)))
				Expect(result.StartStatus).To(Equal(0))
				Expect(result.StopStatus).To(Equal(0))
			})
		})

		Context("When push/start fails", func() {
			It("Calls CF logs", func() {
				fakeCf.StubFailingCF("push")
				sli.RunTest("fake_app_name", "fake_buildpack", "./fake_path", config)

				Expect(fakeCf).To(HaveReceived("RunCF").With(expected_logs_calls))
			})

			It("Cleans up the app", func() {
				fakeCf.StubFailingCF("push")
				sli.RunTest("fake_app_name", "fake_buildpack", "./fake_path", config)

				Expect(fakeCf).To(HaveReceived("RunCF").With(expected_delete_calls))
				Expect(fakeCf).To(HaveReceived("RunCF").With(expected_logout_calls))
			})

			It("Returns an error from CF", func() {
				fakeCf.StubFailingCF("push")
				_, err := sli.RunTest("fake_app_name", "fake_buildpack", "./fake_path", config)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Running CF command failed: push"))
			})

			It("Does not record time and status", func() {
				fakeCf.StubFailingCF("push")
				result, _ := sli.RunTest("fake_app_name", "fake_buildpack", "./fake_path", config)

				Expect(result.StartTime).To(Equal(time.Duration(0)))
				Expect(result.StopTime).To(Equal(time.Duration(0)))
				Expect(result.StartStatus).To(Equal(0))
				Expect(result.StopStatus).To(Equal(0))
			})
		})

		Context("When stop fails", func() {
			It("Calls CF logs", func() {
				fakeCf.StubFailingCF("stop")
				sli.RunTest("fake_app_name", "fake_buildpack", "./fake_path", config)

				Expect(fakeCf).To(HaveReceived("RunCF").With(expected_logs_calls))
			})

			It("Cleans up the app", func() {
				fakeCf.StubFailingCF("stop")
				sli.RunTest("fake_app_name", "fake_buildpack", "./fake_path", config)

				Expect(fakeCf).To(HaveReceived("RunCF").With(expected_delete_calls))
				Expect(fakeCf).To(HaveReceived("RunCF").With(expected_logout_calls))
			})

			It("Returns an error from CF", func() {
				fakeCf.StubFailingCF("stop")
				_, err := sli.RunTest("fake_app_name", "fake_buildpack", "./fake_path", config)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Running CF command failed: stop"))
			})

			It("Records time and status", func() {
				fakeCf.StubFailingCF("stop")
				result, err := sli.RunTest("fake_app_name", "fake_buildpack", "./fake_path", config)
				Expect(err).To(HaveOccurred())

				Expect(result.StartTime).ToNot(Equal(time.Duration(0)))
				Expect(result.StopTime).To(Equal(time.Duration(0)))
				Expect(result.StartStatus).To(Equal(1))
				Expect(result.StopStatus).To(Equal(0))
			})
		})
	})

	Context("#CreateService", func() {
		It("Creates the service", func() {
			err := sli.CreateService("fakeServiceName", "fakePlan", "fakeServiceInstanceName")
			Expect(err).NotTo(HaveOccurred())

			expectedCreateServiceCalls := []string{"create-service", "fakeServiceName", "fakePlan", "fakeServiceInstanceName"}
			Expect(fakeCf).To(HaveReceived("RunCF").With(expectedCreateServiceCalls))
		})

		It("Gets info of the service", func() {
			err := sli.CreateService("fakeServiceName", "fakePlan", "fakeServiceInstanceName")
			Expect(err).NotTo(HaveOccurred())

			expectedServiceCalls := []string{"service", "fakeServiceInstanceName"}
			Expect(fakeCf).To(HaveReceived("RunCF").With(expectedServiceCalls))
		})

		It("Returns error when cf create-service fails", func() {
			fakeCf.StubFailingCF("create-service")
			err := sli.CreateService("fakeServiceName", "fakePlan", "fakeServiceInstanceName")

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Running CF command failed: create-service"))
		})
	})

	Context("#CleanupService", func() {
		It("Deletes the service and logs out", func() {
			err := sli.CleanupService("fakeServiceInstanceName")
			Expect(err).NotTo(HaveOccurred())

			expectedDeleteServiceCalls := []string{"delete-service", "fakeServiceInstanceName", "-f"}
			expectedLogoutCalls := []string{"logout"}
			Expect(fakeCf).To(HaveReceived("RunCF").With(expectedDeleteServiceCalls))
			Expect(fakeCf).To(HaveReceived("RunCF").With(expectedLogoutCalls))
		})

		It("Returns error when cf delete-service fails", func() {
			fakeCf.StubFailingCF("delete-service")
			err := sli.CleanupService("fakeServiceInstanceName")

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Running CF command failed: delete-service"))
		})
	})
})
