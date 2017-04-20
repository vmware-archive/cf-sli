package cf_wrapper_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCfWrapper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CfWrapper Suite")
}
