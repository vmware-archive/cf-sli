package main

import (
	"fmt"
	"os"
	"time"

	"github.com/cloudfoundry-incubator/cf-test-helpers/cf"
	uuid "github.com/nu7hatch/gouuid"
	gexec "github.com/onsi/gomega/gexec"
	"github.com/pivotal-cloudops/cf-sli/config"
)

func wait_for_cf(context *gexec.Session, cf_command string) {
	_ = <-context.Exited
	if context.ExitCode() != 0 {
		fmt.Println("cf " + cf_command + " failed")
		os.Exit(1)
	}
}

func main() {
	var c config.Config
	err := c.LoadConfig("./.config")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	context := cf.Cf("api", c.Api)
	wait_for_cf(context, "api")

	context = cf.Cf("auth", c.User, c.Password)
	wait_for_cf(context, "auth")

	context = cf.Cf("target", "-o", c.Org, "-s", c.Space)
	wait_for_cf(context, "target")

	guid, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	app_name := guid.String()[0:20]

	context = cf.Cf("push", "-p", "./assets/ruby_simple", app_name, "--no-start")
	wait_for_cf(context, "push")

	start := time.Now()
	context = cf.Cf("start", app_name)
	wait_for_cf(context, "start")
	cf_start_elapsed := time.Since(start)

	start = time.Now()
	context = cf.Cf("stop", app_name)
	wait_for_cf(context, "stop")
	cf_stop_elapsed := time.Since(start)

	fmt.Printf("App route: %s.cfapps.io\n", app_name)
	fmt.Printf("App start time: %s\n", cf_start_elapsed)
	fmt.Printf("App stop time: %s\n", cf_stop_elapsed)
}
