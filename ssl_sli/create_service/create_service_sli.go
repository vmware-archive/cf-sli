package createservice

import (
	"bytes"
	"fmt"
	"io"
	"os"

	uuid "github.com/nu7hatch/gouuid"
	"github.com/pivotal-cloudops/cf-sli/config"
	"github.com/pivotal-cloudops/cf-sli/sli_executor"
	datadoghelpers "github.com/pivotal-cloudops/cf-sli/ssl_sli/datadog_helpers"
)

func SLI(config config.Config, sliExecutor sli_executor.SliExecutor, datadogInfo datadoghelpers.DatadogInfo) string {
	fmt.Println("Logging in and targetting space")
	errPrepare := sliExecutor.Prepare(config.Api, config.User, config.Password, config.Org, config.Space)
	if errPrepare != nil {
		fmt.Fprintln(os.Stderr, "There was an error preparing the SLI tool:", errPrepare.Error())
		fmt.Fprintln(os.Stderr, "Exiting and skipping post to Datadog")
		return errPrepare.Error()
	}

	guid, errorGuid := uuid.NewV4()
	if errorGuid != nil {
		panic(errorGuid)
	}
	serviceInstanceName := "ssl-sli-servInst-" + guid.String()[0:18]
	var createStatus int
	err := ""

	fmt.Println("Creating SSL service instance")
	errCreateService := sliExecutor.CreateService("ssl", "free", serviceInstanceName)
	if errCreateService != nil {
		createStatus = 0
		err += errCreateService.Error()
	} else {
		createStatus = 1
	}
	fmt.Println("Create status:", createStatus, "for metric ", datadogInfo.Metric)
	datadoghelpers.PostToDatadog(createStatus, datadogInfo)

	fmt.Println("Cleaning up and deleting the SSL service instance")
	errCleanupService := sliExecutor.CleanupService(serviceInstanceName)
	if errCleanupService != nil {
		err += errCleanupService.Error()
	}

	return err
}

func CaptureStdout(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}
