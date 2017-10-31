package datadoghelpers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestDatadogHelpers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DatadogHelpers Suite")
}
