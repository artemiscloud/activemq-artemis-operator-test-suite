package basic

import (
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/framework/ginkgowrapper"
	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/test"
	"testing"
)

func TestBasic(t *testing.T) {

	gomega.RegisterFailHandler(ginkgowrapper.Fail)
	test.PrepareNamespace(t, "basic", "Basic Suite")
}
