package versions

import (
	"testing"

	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/test"
)

func TestVersions(t *testing.T) {
	test.PrepareNamespace(t, "versions", "Versions Suite")
}

func TestMain(m *testing.M) {
	test.Initialize(m)
}
