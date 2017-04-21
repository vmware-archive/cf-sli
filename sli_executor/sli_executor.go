package sli_executor

import (
	"errors"
	"strings"

	"github.com/pivotal-cloudops/cf-sli/cf_wrapper"
)

type SliExecutor struct {
	Cf_wrapper cf_wrapper.CfWrapperInterface
}

func NewSliExecutor(cf_wrapper cf_wrapper.CfWrapperInterface) *SliExecutor {
	return &SliExecutor{
		Cf_wrapper: cf_wrapper,
	}
}

func (s *SliExecutor) Run(commands ...string) error {
	context := s.Cf_wrapper.RunCF(commands...)
	_ = <-context.Exited
	if context.ExitCode() != 0 {
		error_message := "Running CF command failed: "
		error_message += strings.Join(commands, " ")
		return errors.New(error_message)
	}
	return nil
}
