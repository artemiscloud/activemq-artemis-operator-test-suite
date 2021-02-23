package addresssettings

import (
	"github.com/artemiscloud/activemq-artemis-operator-test-suite/test"
	"github.com/onsi/ginkgo"
)

// Constants available for all test specs related with the One Interior topology
const (
	DeployName = "addrst"
	BaseName   = "brkr"
	//	Command      = "curl"
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
