package main

import (
	"fmt"
	"os"

	"github.com/pivotal-cloudops/cf-sli/cf_wrapper"
	"github.com/pivotal-cloudops/cf-sli/config"
	"github.com/pivotal-cloudops/cf-sli/sli_executor"
	"github.com/pivotal-cloudops/cf-sli/ssl_sli/create_service"
	datadoghelpers "github.com/pivotal-cloudops/cf-sli/ssl_sli/datadog_helpers"
)

func main() {
	var config config.Config
	datadogInfo := datadoghelpers.DatadogInfo{
		DatadogAPIKey:  os.Getenv("PROD_DATADOG_API_KEY"),
		DatadogAppKey:  os.Getenv("PROD_DOGSHELL_DATADOG_APP_KEY"),
		DeploymentName: os.Getenv("BOSH_DEPLOYMENT_NAME"),
	}

	err := config.LoadConfig("./.config")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to load .config:", err)
		os.Exit(1)
	}

	if len(os.Args) == 1 {
		fmt.Fprintln(os.Stderr, "Failed to specify SLI command")
		os.Exit(1)
	}
	sliCommand := os.Args[1]

	var cfCli cf_wrapper.CfWrapper
	sliExecutor := sli_executor.NewSliExecutor(cfCli)

	switch sliCommand {
	case "create-service":
		fmt.Println("Running SSL Service Create Service SLI")
		datadogInfo.Metric = "pws.ssl_service_sli.create_service"
		sliErr := createservice.SLI(config, sliExecutor, datadogInfo)
		if sliErr != "" {
			fmt.Println("Something went wrong when exeucting the SLI: ", sliErr)
			os.Exit(1)
		}
		fmt.Println("Finished running the SSL Service Create Service SLI")
	}
}
