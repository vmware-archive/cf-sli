package sli_executor

import (
	"time"

	"github.com/pivotal-cloudops/cf-sli/cf_wrapper"
	"github.com/pivotal-cloudops/cf-sli/config"
)

type SliExecutor struct {
	Cf_wrapper cf_wrapper.CfWrapperInterface
}

type Result struct {
	StartTime   time.Duration
	StopTime    time.Duration
	StartStatus int
	StopStatus  int
}

func NewSliExecutor(cf_wrapper cf_wrapper.CfWrapperInterface) SliExecutor {
	return SliExecutor{
		Cf_wrapper: cf_wrapper,
	}
}

func (s SliExecutor) cf(commands ...string) error {
	return s.Cf_wrapper.RunCF(commands...)
}

func (s SliExecutor) Prepare(api string, user string, password string, org string, space string) error {
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

func (s SliExecutor) PushAndStartSli(app_name string, app_buildpack string, domain string, path string) (time.Duration, error) {
	err := s.cf("push", "-p", path, "-b", app_buildpack, app_name, "-d", domain, "--no-start")
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

func (s SliExecutor) StopSli(app_name string) (time.Duration, error) {
	start := time.Now()
	err := s.cf("stop", app_name)
	if err != nil {
		return time.Duration(0), err
	}
	time_elapsed := time.Since(start)
	return time_elapsed, nil
}

func (s SliExecutor) CleanupSli(app_name string) error {
	err_delete := s.cf("delete", app_name, "-f")
	err_logout := s.cf("logout")

	if err_delete != nil || err_logout != nil {
		if err_delete != nil {
			return err_delete
		}
		return err_logout
	}

	return nil
}

func (s SliExecutor) RunTest(app_name string, app_buildpack string, path string, c config.Config) (*Result, error) {
	err := s.Prepare(c.Api, c.User, c.Password, c.Org, c.Space)

	if err != nil {
		result := &Result{
			StartStatus: 0,
			StopStatus:  0,
		}
		return result, nil
	}

	defer s.CleanupSli(app_name)
	elapsed_start_time, err := s.PushAndStartSli(app_name, app_buildpack, c.Domain, path)
	if err != nil {
		result := &Result{
			StartStatus: 0,
			StopStatus:  0,
		}
		return result, nil
	}

	elapsed_stop_time, err := s.StopSli(app_name)
	if err != nil {
		result := &Result{
			StartTime:   elapsed_start_time,
			StartStatus: 1,
			StopStatus:  0,
		}
		return result, nil
	}

	result := &Result{
		StartTime:   elapsed_start_time,
		StopTime:    elapsed_stop_time,
		StartStatus: 1,
		StopStatus:  1,
	}
	return result, nil
}
