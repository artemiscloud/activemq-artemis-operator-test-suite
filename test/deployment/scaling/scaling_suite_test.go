package scaling

import (
	"testing"

	"github.com/artemiscloud/activemq-artemis-operator-test-suite/test"
)

func TestBasic(t *testing.T) {
	test.PrepareNamespace(t, "scaling", "Scaling Suite")
}

func TestMain(m *testing.M) {
	test.Initialize(m)
}
