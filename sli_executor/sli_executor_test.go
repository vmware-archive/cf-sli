package sli_executor_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cloudops/cf-sli/cf_wrapper/cf_wrapperfakes"
	"github.com/pivotal-cloudops/cf-sli/config"
	"github.com/pivotal-cloudops/cf-sli/sli_executor"
)

var _ = Describe("SliExecutor", func() {
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
			expected_api_calls := []string{"api", "fake_api"}
			Expect(fakeCf.RunCFArgsForCall(0)).To(Equal(expected_api_calls))
			expected_auth_calls := []string{"auth", "fake_user", "fake_pass"}
			Expect(fakeCf.RunCFArgsForCall(1)).To(Equal(expected_auth_calls))
			expected_target_calls := []string{"target", "-o", "fake_org", "-s", "fake_space"}
			Expect(fakeCf.RunCFArgsForCall(2)).To(Equal(expected_target_calls))
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
			elapsed_time, err := sli.PushAndStartSli("fake_app_name", "fake_buildpack", "fake_domain", "./fake_path")
			Expect(err).NotTo(HaveOccurred())
			expected_push_calls := []string{"push", "-p", "./fake_path", "-b", "fake_buildpack", "fake_app_name", "-d", "fake_domain", "--no-start"}
			Expect(fakeCf.RunCFArgsForCall(0)).To(Equal(expected_push_calls))
			expected_start_calls := []string{"start", "fake_app_name"}
			Expect(fakeCf.RunCFArgsForCall(1)).To(Equal(expected_start_calls))
			Expect(elapsed_time).ToNot(Equal(0))
		})

		It("Returns error when cf push fails", func() {
			fakeCf.StubFailingCF("push")
			elapsed_time, err := sli.PushAndStartSli("fake_app_name", "fake_buildpack", "fake_domain", "./fake_path")
			Expect(err).To(HaveOccurred())
			Expect(elapsed_time).To(Equal(time.Duration(0)))
		})

		It("Returns error when cf start fails", func() {
			fakeCf.StubFailingCF("start")
			elapsed_time, err := sli.PushAndStartSli("fake_app_name", "fake_buildpack", "fake_domain", "./fake_path")
			Expect(err).To(HaveOccurred())
			Expect(elapsed_time).To(Equal(time.Duration(0)))
		})
	})

	Context("#StopSli", func() {
		It("Start the Sli app", func() {
			elapsed_time, err := sli.StopSli("fake_app_name")
			Expect(err).NotTo(HaveOccurred())
			expected_stop_calls := []string{"stop", "fake_app_name"}
			Expect(fakeCf.RunCFArgsForCall(0)).To(Equal(expected_stop_calls))
			Expect(elapsed_time).ToNot(Equal(time.Duration(0)))
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
			expected_delete_calls := []string{"delete", "fake_app_name", "-f"}
			Expect(fakeCf.RunCFArgsForCall(0)).To(Equal(expected_delete_calls))
			expected_logout_calls := []string{"logout"}
			Expect(fakeCf.RunCFArgsForCall(1)).To(Equal(expected_logout_calls))
		})

		It("Returns error when cf delete fails, and it logs out", func() {
			fakeCf.StubFailingCF("delete")
			err := sli.CleanupSli("fake_app_name")
			Expect(err).To(HaveOccurred())
			expected_logout_calls := []string{"logout"}
			Expect(fakeCf.RunCFArgsForCall(1)).To(Equal(expected_logout_calls))
		})

		It("Returns error when cf logout fails", func() {
			fakeCf.StubFailingCF("logout")
			err := sli.CleanupSli("fake_app_name")
			Expect(err).To(HaveOccurred())
		})
	})

	Context("#RunTest", func() {
		It("Login, push the app, returns the start and stop times and status, and cleanup", func() {
			result, err := sli.RunTest("fake_app_name", "fake_buildpack", "./fake_app_path", config)
			Expect(err).NotTo(HaveOccurred())

			// Login and target to the org and space
			expected_api_calls := []string{"api", "fake_api"}
			Expect(fakeCf.RunCFArgsForCall(0)).To(Equal(expected_api_calls))
			expected_auth_calls := []string{"auth", "fake_user", "fake_pass"}
			Expect(fakeCf.RunCFArgsForCall(1)).To(Equal(expected_auth_calls))
			expected_target_calls := []string{"target", "-o", "fake_org", "-s", "fake_space"}
			Expect(fakeCf.RunCFArgsForCall(2)).To(Equal(expected_target_calls))

			// Push, start, and stop the app
			expected_push_calls := []string{"push", "-p", "./fake_app_path", "-b", "fake_buildpack", "fake_app_name", "-d", "fake_domain", "--no-start"}
			Expect(fakeCf.RunCFArgsForCall(3)).To(Equal(expected_push_calls))
			expected_start_calls := []string{"start", "fake_app_name"}
			Expect(fakeCf.RunCFArgsForCall(4)).To(Equal(expected_start_calls))
			expected_stop_calls := []string{"stop", "fake_app_name"}
			Expect(fakeCf.RunCFArgsForCall(5)).To(Equal(expected_stop_calls))
			Expect(result.StartTime).ToNot(Equal(0))
			Expect(result.StopTime).ToNot(Equal(0))
			Expect(result.StartStatus).To(Equal(1))
			Expect(result.StopStatus).To(Equal(1))

			// Cleanup and logout
			expected_delete_calls := []string{"delete", "fake_app_name", "-f"}
			Expect(fakeCf.RunCFArgsForCall(6)).To(Equal(expected_delete_calls))
			expected_logout_calls := []string{"logout"}
			Expect(fakeCf.RunCFArgsForCall(7)).To(Equal(expected_logout_calls))
		})

		It("Cleans up the app if push fails", func() {

			fakeCf.StubFailingCF("push")
			result, err := sli.RunTest("fake_app_name", "fake_buildpack", "./fake_app_path", config)
			Expect(err).NotTo(HaveOccurred())

			// Login and target to the org and space
			expected_api_calls := []string{"api", "fake_api"}
			Expect(fakeCf.RunCFArgsForCall(0)).To(Equal(expected_api_calls))
			expected_auth_calls := []string{"auth", "fake_user", "fake_pass"}
			Expect(fakeCf.RunCFArgsForCall(1)).To(Equal(expected_auth_calls))
			expected_target_calls := []string{"target", "-o", "fake_org", "-s", "fake_space"}
			Expect(fakeCf.RunCFArgsForCall(2)).To(Equal(expected_target_calls))

			// call #3: Push the app will fail

			// Cleanup and logout
			expected_delete_calls := []string{"delete", "fake_app_name", "-f"}
			Expect(fakeCf.RunCFArgsForCall(4)).To(Equal(expected_delete_calls))
			expected_logout_calls := []string{"logout"}
			Expect(fakeCf.RunCFArgsForCall(5)).To(Equal(expected_logout_calls))

			Expect(result.StartTime).To(Equal(time.Duration(0)))
			Expect(result.StopTime).To(Equal(time.Duration(0)))
			Expect(result.StartStatus).To(Equal(0))
			Expect(result.StopStatus).To(Equal(0))
		})

		It("Cleans up the app if stop fails", func() {
			fakeCf.StubFailingCF("stop")
			result, err := sli.RunTest("fake_app_name", "fake_buildpack", "./fake_app_path", config)
			Expect(err).NotTo(HaveOccurred())

			// Login and target to the org and space
			expected_api_calls := []string{"api", "fake_api"}
			Expect(fakeCf.RunCFArgsForCall(0)).To(Equal(expected_api_calls))
			expected_auth_calls := []string{"auth", "fake_user", "fake_pass"}
			Expect(fakeCf.RunCFArgsForCall(1)).To(Equal(expected_auth_calls))
			expected_target_calls := []string{"target", "-o", "fake_org", "-s", "fake_space"}
			Expect(fakeCf.RunCFArgsForCall(2)).To(Equal(expected_target_calls))

			// Push and start app
			expected_push_calls := []string{"push", "-p", "./fake_app_path", "-b", "fake_buildpack", "fake_app_name", "-d", "fake_domain", "--no-start"}
			Expect(fakeCf.RunCFArgsForCall(3)).To(Equal(expected_push_calls))
			expected_start_calls := []string{"start", "fake_app_name"}
			Expect(fakeCf.RunCFArgsForCall(4)).To(Equal(expected_start_calls))
			// call #5: stop the app will fail

			// Cleanup and logout
			expected_delete_calls := []string{"delete", "fake_app_name", "-f"}
			Expect(fakeCf.RunCFArgsForCall(6)).To(Equal(expected_delete_calls))
			expected_logout_calls := []string{"logout"}
			Expect(fakeCf.RunCFArgsForCall(7)).To(Equal(expected_logout_calls))

			Expect(result.StartTime).ToNot(Equal(time.Duration(0)))
			Expect(result.StopTime).To(Equal(time.Duration(0)))
			Expect(result.StartStatus).To(Equal(1))
			Expect(result.StopStatus).To(Equal(0))
		})
	})
})
