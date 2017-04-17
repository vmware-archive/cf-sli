package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/cloudfoundry-incubator/cf-test-helpers/cf"
	uuid "github.com/nu7hatch/gouuid"
	gexec "github.com/onsi/gomega/gexec"
	"github.com/pivotal-cloudops/cf-sli/config"
)

type Result struct {
	Route     string `json:"app_route"`
	StartTime string `json:"app_start_time"`
	StopTime  string `json:"app_stop_time"`
}

func wait_for_cf(context *gexec.Session, cf_command string) {
	_ = <-context.Exited
	if context.ExitCode() != 0 {
		fmt.Println("cf " + cf_command + " failed")
		os.Exit(1)
	}
}

func clean_up(app_name string) {
	context := cf.Cf("delete", app_name, "-f")
	wait_for_cf(context, "delete")
	context = cf.Cf("logout")
	wait_for_cf(context, "logout")
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

	defer clean_up(app_name)
	context = cf.Cf("push", "-p", "./assets/ruby_simple", app_name, "-d", c.Domain, "--no-start")
	wait_for_cf(context, "push")

	start := time.Now()
	context = cf.Cf("start", app_name)
	wait_for_cf(context, "start")
	cf_start_elapsed := time.Since(start)

	start = time.Now()
	context = cf.Cf("stop", app_name)
	wait_for_cf(context, "stop")
	cf_stop_elapsed := time.Since(start)

	result := &Result{
		Route:     app_name + "." + c.Domain,
		StartTime: cf_start_elapsed.String(),
		StopTime:  cf_stop_elapsed.String(),
	}

	output, _ := json.Marshal(result)
	fmt.Println(string(output))
}
