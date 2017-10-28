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

	fmt.Println("Creating SSL service instance")
	guid := generateGuid()
	serviceInstanceName := "ssl-sli-servInst-" + guid[0:18]

	createStatus, err := performSLITask(serviceInstanceName, sliExecutor)
	if err != "" {
		fmt.Println("An error occured when executing the SLI: ", err)
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

func performSLITask(serviceInstanceName string, sliExecutor sli_executor.SliExecutor) (int, string) {
	err := sliExecutor.CreateService("ssl", "free", serviceInstanceName)
	if err != nil {
		return 0, err.Error()
	}
	return 1, ""
}

func generateGuid() string {
	guid, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	return guid.String()
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
