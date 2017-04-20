package cf_wrapper

import (
	"errors"
	"strings"

	"github.com/cloudfoundry-incubator/cf-test-helpers/cf"
)

type CfWrapperInterface interface {
	RunCF(commands ...string) error
}

type CfWrapper struct {
}

func (h *CfWrapper) RunCF(commands ...string) error {
	context := cf.Cf(commands...)
	_ = <-context.Exited
	if context.ExitCode() != 0 {
		error_message := "Running CF command failed: "
		error_message += strings.Join(commands, " ")
		return errors.New(error_message)
	}
	return nil
}
