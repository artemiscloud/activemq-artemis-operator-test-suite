package routes

import (
	"testing"

	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/test"
)

func TestBasic(t *testing.T) {
	test.PrepareNamespace(t, "routes", "Routes Suite")
}

func TestMain(m *testing.M) {
	test.Initialize(m)
}
