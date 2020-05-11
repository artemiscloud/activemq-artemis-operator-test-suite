package basic

import (
	"testing"

	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/test"
)

func TestBasic(t *testing.T) {
	//gomega.RegisterFailHandler(ginkgowrapper.Fail)
	test.PrepareNamespace(t, "basic", "Basic Suite")
}

func TestMain(m *testing.M) {
	test.Initialize(m)
}
