package statistics

import (
	"testing"

	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/test"
)

func TestBasic(t *testing.T) {
	test.PrepareNamespace(t, "statistics", "Statistics Suite")
}

func TestMain(m *testing.M) {
	test.Initialize(m)
}
