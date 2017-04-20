package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	uuid "github.com/nu7hatch/gouuid"
	"github.com/pivotal-cloudops/cf-sli/cf_wrapper"
	"github.com/pivotal-cloudops/cf-sli/config"
)

type Result struct {
	Route     string `json:"app_route"`
	StartTime string `json:"app_start_time"`
	StopTime  string `json:"app_stop_time"`
}

func clean_up(app_name string, wrapper cf_wrapper.CfWrapperInterface) {
	err := wrapper.RunCF("delete", app_name, "-f")
	if err != nil {
		panic(err)
	}
	err = wrapper.RunCF("logout")
	if err != nil {
		panic(err)
	}
}

func main() {
	var c config.Config
	var wrapper *cf_wrapper.CfWrapper

	err := c.LoadConfig("./.config")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = wrapper.RunCF("api", c.Api)
	if err != nil {
		panic(err)
	}

	err = wrapper.RunCF("auth", c.User, c.Password)
	if err != nil {
		panic(err)
	}

	err = wrapper.RunCF("target", "-o", c.Org, "-s", c.Space)
	if err != nil {
		panic(err)
	}

	guid, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	app_name := guid.String()[0:20]

	defer clean_up(app_name, wrapper)
	err = wrapper.RunCF("push", "-p", "./assets/ruby_simple", app_name, "-d", c.Domain, "--no-start")
	if err != nil {
		panic(err)
	}

	start := time.Now()
	err = wrapper.RunCF("start", app_name)
	if err != nil {
		panic(err)
	}

	cf_start_elapsed := time.Since(start)

	start = time.Now()
	err = wrapper.RunCF("stop", app_name)
	if err != nil {
		panic(err)
	}

	cf_stop_elapsed := time.Since(start)

	result := &Result{
		Route:     app_name + "." + c.Domain,
		StartTime: cf_start_elapsed.String(),
		StopTime:  cf_stop_elapsed.String(),
	}

	output, _ := json.Marshal(result)
	fmt.Fprintf(os.Stderr, string(output))
}
