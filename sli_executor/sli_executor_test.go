package sli_executor_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	"github.com/pivotal-cloudops/cf-sli/cf_wrapper/cf_wrapperfakes"
	"github.com/pivotal-cloudops/cf-sli/sli_executor"
)

var _ = Describe("Run CF", func() {
	var (
		fakeCf *cf_wrapperfakes.FakeCfWrapperInterface
	)

	BeforeEach(func() {
		fakeCf = new(cf_wrapperfakes.FakeCfWrapperInterface)
	})

	It("returns nil if cf command executes successfully", func() {
		sli := sli_executor.NewSliExecutor(fakeCf)
		sli.Run("help")
		fmt.Println(fakeCf.RunCFArgsForCall(0))
	})
})
