package sli_executor

import (
	"time"

	"os"
	"strconv"
	"strings"

	"github.com/pivotal-cloudops/cf-sli/cf_wrapper"
	"github.com/pivotal-cloudops/cf-sli/config"
	"github.com/pivotal-cloudops/cf-sli/logger"
)

type SliExecutor struct {
	Cf_wrapper cf_wrapper.CfWrapperInterface
	logger     logger.Logger
}

type Result struct {
	StartTime   time.Duration
	StopTime    time.Duration
	StartStatus int
	StopStatus  int
}

func NewSliExecutor(cf_wrapper cf_wrapper.CfWrapperInterface, logger logger.Logger) SliExecutor {
	return SliExecutor{
		Cf_wrapper: cf_wrapper,
		logger:     logger,
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

func (s SliExecutor) PushAndStartSli(app_name string, domain string, path string, timeouts config.TimeoutConfig) (time.Duration, error) {

	s.logger.Printf("PUSH_TIMEOUTS: %+v", timeouts)

	os.Setenv("CF_STAGING_TIMEOUT", strconv.Itoa(timeouts.Staging))
	os.Setenv("CF_STARTUP_TIMEOUT", strconv.Itoa(timeouts.Startup))

	manifest := strings.Join([]string{path, "/manifest.yml"}, "")

	err := s.cf("push", "-p", path, app_name, "-f", manifest, "-d", domain, "--no-start", "-t", strconv.Itoa(timeouts.FirstHealthyResponse))
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
	err_delete := s.cf("delete", app_name, "-f", "-r")
	err_logout := s.cf("logout")

	if err_delete != nil || err_logout != nil {
		if err_delete != nil {
			return err_delete
		}
		return err_logout
	}

	return nil
}

func (s SliExecutor) RunTest(app_name string, path string, config config.Config) (*Result, error) {
	defer s.CleanupSli(app_name)

	err := s.Prepare(config.Api, config.User, config.Password, config.Org, config.Space)
	if err != nil {
		result := &Result{
			StartStatus: 0,
			StopStatus:  0,
		}
		return result, err
	}

	elapsedStartTime, err := s.PushAndStartSli(app_name, config.Domain, path, config.Timeout)
	if err != nil {
		result := &Result{
			StartStatus: 0,
			StopStatus:  0,
		}
		s.printLogs(app_name)
		return result, err
	}

	elapsedStopTime, err := s.StopSli(app_name)
	if err != nil {
		result := &Result{
			StartTime:   elapsedStartTime,
			StartStatus: 1,
			StopStatus:  0,
		}
		s.printLogs(app_name)
		return result, err
	}

	result := &Result{
		StartTime:   elapsedStartTime,
		StopTime:    elapsedStopTime,
		StartStatus: 1,
		StopStatus:  1,
	}
	return result, nil
}

func (s SliExecutor) printLogs(app_name string) {
	s.cf("app", app_name, "--guid")
	s.cf("logs", app_name, "--recent")
}

func (s SliExecutor) CreateService(serviceName string, plan string, serviceInstanceName string) error {
	err := s.cf("create-service", serviceName, plan, serviceInstanceName)
	if err != nil {
		return err
	}

	err = s.cf("service", serviceInstanceName)
	if err != nil {
		return err
	}

	return nil
}

func (s SliExecutor) CleanupService(serviceInstanceName string) error {
	err := s.cf("delete-service", serviceInstanceName, "-f")
	if err != nil {
		return err
	}

	err = s.cf("logout")
	if err != nil {
		return err
	}

	return nil
}
