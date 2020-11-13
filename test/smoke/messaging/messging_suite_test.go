package messaging

import (
	"testing"

	"github.com/artemiscloud/activemq-artemis-operator-test-suite/test"
)

func TestMessaging(t *testing.T) {
	test.PrepareNamespace(t, "messaging", "Messaging Suite")
}

func TestMain(m *testing.M) {
	test.Initialize(m)
}
