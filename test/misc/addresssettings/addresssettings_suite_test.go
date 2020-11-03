package addresssettings

import (
	"testing"

	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/test"
)

func TestMessaging(t *testing.T) {
	test.PrepareNamespace(t, "addresssetting", "Address Setting test suite")
}

func TestMain(m *testing.M) {
	test.Initialize(m)
}
