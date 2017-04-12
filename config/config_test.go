package config_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cloudops/cf-sli/config"
)

var _ = Describe("LoadConfig", func() {
	It("loads config from a JSON file", func() {
		var c config.Config
		err := c.LoadConfig("../fixtures/config_test.json")
		Expect(err).ToNot(HaveOccurred())
		Expect(c.Api).To(Equal("fake_api"))
		Expect(c.User).To(Equal("fake_user"))
		Expect(c.Org).To(Equal("fake_org"))
		Expect(c.Space).To(Equal("fake_space"))
		Expect(c.Password).To(Equal("fake_pass"))
	})

	It("returns an error reading a none-existing file", func() {
		var c config.Config
		err := c.LoadConfig("../fixtures/none_existing_config_test.json")
		Expect(err).To(HaveOccurred())
	})
})
