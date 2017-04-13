package main

import (
	"fmt"
	"os"

	"github.com/cloudfoundry-incubator/cf-test-helpers/cf"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/pivotal-cloudops/cf-sli/config"
)

func main() {
	var c config.Config
	err := c.LoadConfig("./.config")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	context := cf.Cf("api", c.Api)
	_ = <-context.Exited

	if context.ExitCode() != 0 {
		fmt.Println("cf api failed")
		os.Exit(1)
	}

	context = cf.Cf("auth", c.User, c.Password)
	_ = <-context.Exited

	if context.ExitCode() != 0 {
		fmt.Println("cf auth failed")
		os.Exit(1)
	}

	context = cf.Cf("target", "-o", c.Org, "-s", c.Space)
	_ = <-context.Exited

	if context.ExitCode() != 0 {
		fmt.Println("cf target failed")
		os.Exit(1)
	}

	guid, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	context = cf.Cf("push", "-p", "./assets/ruby_simple", guid.String()[0:20], "--no-start")
	_ = <-context.Exited

	if context.ExitCode() != 0 {
		fmt.Println("cf push failed")
		os.Exit(1)
	}

}
