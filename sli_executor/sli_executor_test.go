package sli_executor_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cloudops/cf-sli/cf_wrapper/cf_wrapperfakes"
	"github.com/pivotal-cloudops/cf-sli/sli_executor"
)

var _ = Describe("SliExecutor", func() {
	var (
		fakeCf *cf_wrapperfakes.FakeCfWrapperInterface
		sli    *sli_executor.SliExecutor
	)
	BeforeEach(func() {
		fakeCf = new(cf_wrapperfakes.FakeCfWrapperInterface)
		sli = sli_executor.NewSliExecutor(fakeCf)
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
		It("Push the Sli app with --no-start", func() {
			_, err := sli.PushAndStartSli("fake_app_name", "fake_domain", "./fake_path")
			Expect(err).NotTo(HaveOccurred())
			expected_push_calls := []string{"push", "-p", "./fake_path", "fake_app_name", "-d", "fake_domain", "--no-start"}
			Expect(fakeCf.RunCFArgsForCall(0)).To(Equal(expected_push_calls))
		})

		It("Runs cf start and returns how long it takes", func() {
			elapsed_time, err := sli.PushAndStartSli("fake_app_name", "fake_domain", "./fake_path")
			Expect(err).NotTo(HaveOccurred())
			expected_start_calls := []string{"start", "fake_app_name"}
			Expect(fakeCf.RunCFArgsForCall(1)).To(Equal(expected_start_calls))
			Expect(elapsed_time).ToNot(Equal(0))
		})

		It("Returns error when cf push fails", func() {
			fakeCf.StubFailingCF("push")
			elapsed_time, err := sli.PushAndStartSli("fake_app_name", "fake_domain", "./fake_path")
			Expect(err).To(HaveOccurred())
			Expect(elapsed_time).To(Equal(time.Duration(0)))
		})

		It("Returns error when cf start fails", func() {
			fakeCf.StubFailingCF("start")
			elapsed_time, err := sli.PushAndStartSli("fake_app_name", "fake_domain", "./fake_path")
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
})
