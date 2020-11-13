package addresssettings

import (
	"testing"

	"github.com/artemiscloud/activemq-artemis-operator-test-suite/test"
)

func TestAddressSettings(t *testing.T) {
	test.PrepareNamespace(t, "addresssetting", "Address Setting test suite")
}

func TestMain(m *testing.M) {
	test.Initialize(m)
}
