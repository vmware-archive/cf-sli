package cf_wrapper

import (
	"github.com/cloudfoundry-incubator/cf-test-helpers/cf"
	"github.com/onsi/gomega/gexec"
)

type CfWrapperInterface interface {
	RunCF(commands ...string) *gexec.Session
}

type CfWrapper struct {
}

func (h *CfWrapper) RunCF(commands ...string) *gexec.Session {
	return cf.Cf(commands...)
}
