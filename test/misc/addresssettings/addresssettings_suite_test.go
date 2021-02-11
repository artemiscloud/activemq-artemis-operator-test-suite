package addresssettings

import (
	"github.com/artemiscloud/activemq-artemis-operator-test-suite/pkg/bdw"
	"github.com/artemiscloud/activemq-artemis-operator-test-suite/test"
	"github.com/rh-messaging/shipshape/pkg/framework"
	"testing"
   	"github.com/rh-messaging/shipshape/pkg/framework/log"
	"encoding/json"
   	"github.com/artemiscloud/activemq-artemis-operator-test-suite/pkg/test_helpers"
	"strings"
)

func TestAddressSettings(t *testing.T) {
	test.PrepareNamespace(t, "addresssetting", "Address Setting test suite")
}

func TestMain(m *testing.M) {
	test.Initialize(m)
}

func setEnv(ctx1 *framework.ContextData, brokerDeployer *bdw.BrokerDeploymentWrapper) {
	brokerDeployer.WithWait(true).
		WithContext(ctx1).
		WithBrokerClient(sw.BrokerClient).
		WithCustomImage(test.Config.BrokerImageName).
		WithName(DeployName).
		WithLts(!test.Config.NeedsLatestCR).
		WithIncreasedTimeout(test.Config.TimeoutMultiplier).
        WithConsoleExposure(true)
}

func retrieveAddressSettings(address, AddressBit string, hw *test_helpers.HttpWrapper) test_helpers.Value {
	domain := strings.Split(address, ".")[0]
	header := strings.Replace(OriginHeader, "NAME", domain, 1)
	hw.AddHeader("Origin", header)
	actualURL := "http://" + address + CallAddress + AddressBit
	log.Logf("url (string): %s", actualURL)
	hw.WithPassword(test.Username).WithUser(test.Password)
	result, err := hw.PerformHttpRequest(actualURL)
	if err != nil {
		log.Logf("%s", err)
	}
	//result = strings.ReplaceAll(result,"\\\"","\"")
	var item test_helpers.JolokiaBrokerSettings
	json.Unmarshal([]byte(result), &item)
	var value test_helpers.Value
	json.Unmarshal([]byte(item.Value), &value)
	return value
}
