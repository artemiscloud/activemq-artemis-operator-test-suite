package configuration

import (
	"testing"

	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/test"
)

func TestConfiguration(t *testing.T) {
	test.PrepareNamespace(t, "configuration", "Configuration Suite")
}

func TestMain(m *testing.M) {
	test.Initialize(m)
}
