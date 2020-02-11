package basic_test

import (
	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/test"
	"github.com/rh-messaging/shipshape/pkg/framework/ginkgowrapper"
	"testing"
	 "github.com/onsi/ginkgo"
	 "github.com/onsi/gomega"
)

func TestBasic(t *testing.T) {

	gomega.RegisterFailHandler(ginkgowrapper.Fail)
	test.Initialize(t, "basic", "Basic Suite")

	ginkgo.RunSpecs(t, "Basic Suite")

}
