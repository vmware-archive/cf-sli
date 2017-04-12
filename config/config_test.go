package config_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cloudops/cf-sli/config"
)

var _ = Describe("LoadConfig", func() {
	It("loads config from a JSON file", func() {
		var c config.Config
		c.LoadConfig("../fixtures/config_test.json")
		Expect(c.Api).To(Equal("fake_api"))
	})
})
