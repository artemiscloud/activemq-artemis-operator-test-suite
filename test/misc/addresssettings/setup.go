package addresssettings

import (
	"encoding/json"
	"github.com/artemiscloud/activemq-artemis-operator-test-suite/pkg/test_helpers"
	"github.com/artemiscloud/activemq-artemis-operator-test-suite/test"
	"github.com/onsi/ginkgo"
	"github.com/rh-messaging/shipshape/pkg/framework/log"
	"strings"
)

// Constants available for all test specs related with the One Interior topology
const (
	DeployName   = "addrst"
	BaseName     = "brkr"
	Command      = "curl"
	OriginHeader = "http://localhost:8161"
	// addrst-wconsj-0-svc-rte-e2e-tests-brkr-rx9ls.apps.brokerteam

	//             /console/jolokia/exec/org.apache.activemq.artemis:broker=\"0.0.0.0\"/getAddressSettingsAsJSON/DLQ
	CallAddress = "/console/jolokia/exec/org.apache.activemq.artemis:broker=\"amq-broker\"/getAddressSettingsAsJSON/"
)

var (
	sw *test.SetupWrapper
)

// Create the Framework instance to be used oneinterior test

var _ = ginkgo.BeforeEach(func() {
	sw = &test.SetupWrapper{}
	sw.WithBaseName(BaseName).WithDeployName(DeployName)
	if !test.Config.NeedsLatestCR {
		ginkgo.Skip("Not supported on pre-0.17 operator")
	}
	sw.BeforeEach()
}, 60)

var _ = ginkgo.JustBeforeEach(func() {
	sw.JustBeforeEach()
})

// After each test completes, run cleanup actions to save resources (otherwise resources will remain till
// all specs from this suite are done.
var _ = ginkgo.AfterEach(func() {
	sw.AfterEach()
})

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
