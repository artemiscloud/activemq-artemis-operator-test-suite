package addresssettings

import (
    "strings"
    "strconv"
    "encoding/json"
	"github.com/onsi/ginkgo"
    "github.com/onsi/gomega"
	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/test"
    "github.com/rh-messaging/shipshape/pkg/framework/log"
    "gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/pkg/test_helpers"
)

// Constants available for all test specs related with the One Interior topology
const (
	DeployName = "addrst"
	BaseName   = "brkr"

	Command = "curl"

	OriginHeader = "http://NAME.svc.cluster.local"
	// addrst-wconsj-0-svc-rte-e2e-tests-brkr-rx9ls.apps.brokerteam
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

func verifyAddressSettingsInt(address, AddressBit, itemName string,
                              value int, hw *test_helpers.HttpWrapper) {
	domain := strings.Split(address, ".")[0]
	header := strings.Replace(OriginHeader, "NAME", domain, 1)
	hw.AddHeader("Origin", header)
	actualUrl := "http://admin:admin@" + address + CallAddress + AddressBit
	hw.WithPassword("admin").WithUser("admin")
	result, err := hw.PerformHttpRequest(actualUrl)
	if err != nil {
		log.Logf("%s", err)
	}
	var item map[string]map[string]string
	json.Unmarshal([]byte(result), &item)
	brokerValue := item["value"][itemName]
	gomega.Expect(brokerValue).To(gomega.Equal(string(value)))

}

func verifyAddressSettingsBool(address, AddressBit, itemName string,
                              value bool, hw *test_helpers.HttpWrapper) {
	domain := strings.Split(address, ".")[0]
	header := strings.Replace(OriginHeader, "NAME", domain, 1)
	hw.AddHeader("Origin", header)
	actualUrl := "http://admin:admin@" + address + CallAddress + AddressBit
	hw.WithPassword("admin").WithUser("admin")
	result, err := hw.PerformHttpRequest(actualUrl)
	if err != nil {
		log.Logf("%s", err)
	}
	var item map[string]map[string]string
	json.Unmarshal([]byte(result), &item)
	brokerValue := item["value"][itemName]
	gomega.Expect(strconv.ParseBool(brokerValue)).To(gomega.Equal(value))

}

func verifyAddressSettingsString(address, AddressBit, itemName string,
                              value string, hw *test_helpers.HttpWrapper) {
	domain := strings.Split(address, ".")[0]
	header := strings.Replace(OriginHeader, "NAME", domain, 1)
	hw.AddHeader("Origin", header)
	actualUrl := "http://admin:admin@" + address + CallAddress + AddressBit
	hw.WithPassword("admin").WithUser("admin")
	result, err := hw.PerformHttpRequest(actualUrl)
	if err != nil {
		log.Logf("%s", err)
	}
	var item map[string]map[string]string
	json.Unmarshal([]byte(result), &item)
	brokerValue := item["value"][itemName]
	gomega.Expect(brokerValue).To(gomega.Equal(value))

}

func verifyAddressSettingsFloat(address, AddressBit, itemName string,
                              value float64, hw *test_helpers.HttpWrapper) {
	domain := strings.Split(address, ".")[0]
	header := strings.Replace(OriginHeader, "NAME", domain, 1)
	hw.AddHeader("Origin", header)
	actualUrl := "http://admin:admin@" + address + CallAddress + AddressBit
	hw.WithPassword("admin").WithUser("admin")
	result, err := hw.PerformHttpRequest(actualUrl)
	if err != nil {
		log.Logf("%s", err)
	}
	var item map[string]map[string]string
	json.Unmarshal([]byte(result), &item)
	brokerValue := item["value"][itemName]
	gomega.Expect(brokerValue).To(gomega.Equal(strconv.FormatFloat(value, 'f', 1, 64)))

}
