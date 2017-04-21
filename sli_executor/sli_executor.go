package sli_executor

import (
	"time"

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

func (s *SliExecutor) cf(commands ...string) error {
	return s.Cf_wrapper.RunCF(commands...)
}

func (s *SliExecutor) Prepare(api string, user string, password string, org string, space string) error {
	err := s.cf("api", api)
	if err != nil {
		return err
	}
	err = s.cf("auth", user, password)
	if err != nil {
		return err
	}
	err = s.cf("target", "-o", org, "-s", space)
	if err != nil {
		return err
	}
	return nil
}

func (s *SliExecutor) PushSli(app_name string, domain string, path string) (time.Duration, error) {
	err := s.cf("push", "-p", path, app_name, "-d", domain, "--no-start")
	if err != nil {
		return time.Duration(0), err
	}

	start := time.Now()
	err = s.cf("start", app_name)
	if err != nil {
		return time.Duration(0), err
	}

	time_elapsed := time.Since(start)
	return time_elapsed, nil
}
