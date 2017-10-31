package createservice

import (
	"bytes"
	"log"
	"os"

	uuid "github.com/nu7hatch/gouuid"
	"github.com/pivotal-cloudops/cf-sli/config"
	"github.com/pivotal-cloudops/cf-sli/sli_executor"
	datadoghelpers "github.com/pivotal-cloudops/cf-sli/ssl_sli/datadog_helpers"
)

func SLI(config config.Config, sliExecutor sli_executor.SliExecutor, datadogInfo datadoghelpers.DatadogInfo) string {
	err := setup(config, sliExecutor)
	if err != nil {
		log.Println("Exiting and skipping post to Datadog")
		return err.Error()
	}

	log.Println("Creating SSL service instance")
	guid := generateGuid()
	serviceInstanceName := "ssl-sli-servInst-" + guid[0:18]

	createStatus := 1
	errSLI := ""
	err = sliExecutor.CreateService("ssl", "free", serviceInstanceName)
	if err != nil {
		log.Println("An error occurred when executing the SLI: ", err)
		createStatus = 0
		errSLI = err.Error()
	}
	log.Println("Create status:", createStatus, "for metric ", datadogInfo.Metric)

	respStatus := datadoghelpers.PostToDatadog(createStatus, datadogInfo)
	if respStatus != "200" {
		log.Println("An error occurred when posting to Datadog, response status: ", respStatus)
		errSLI += "Error reporting to Datadog"
	}

	err = cleanup(serviceInstanceName, sliExecutor)
	if err != nil {
		log.Println("An error occurred when cleaning up service instance: ", err)
		errSLI += err.Error()
	}

	return errSLI
}

func setup(config config.Config, sliExecutor sli_executor.SliExecutor) error {
	log.Println("Logging in and targetting space")
	err := sliExecutor.Prepare(config.Api, config.User, config.Password, config.Org, config.Space)
	if err != nil {
		log.Println("There was an error preparing the SLI tool:", err.Error())
		return err
	}
	return nil
}

func cleanup(serviceInstanceName string, sliExecutor sli_executor.SliExecutor) error {
	log.Println("Cleaning up and deleting the SSL service instance")
	err := sliExecutor.CleanupService(serviceInstanceName)
	if err != nil {
		return err
	}
	return nil
}

func generateGuid() string {
	guid, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	return guid.String()
}

func CaptureLog(f func()) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stdout)
	}()

	f()

	return buf.String()
}
