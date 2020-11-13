package scaling

import (
	"testing"

	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/test"
)

func TestBasic(t *testing.T) {
	test.PrepareNamespace(t, "basic", "Basic Suite")
}

func TestMain(m *testing.M) {
	test.Initialize(m)
}
