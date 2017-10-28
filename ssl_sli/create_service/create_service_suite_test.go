package createservice_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCreateService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CreateService Suite")
}
