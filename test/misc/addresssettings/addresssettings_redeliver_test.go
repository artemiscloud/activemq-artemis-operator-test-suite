package addresssettings

import (
	"fmt"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/framework"
	"github.com/rh-messaging/shipshape/pkg/framework/log"
	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/pkg/bdw"
	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/pkg/test_helpers"
	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/test"
)

var _ = ginkgo.Describe("AddressSettingsRedeliveryTest", func() {

	var (
		ctx1 *framework.ContextData
		//brokerClient brokerclientset.Interface
		brokerDeployer *bdw.BrokerDeploymentWrapper
		//url      string
	)

	var (
		AddressBit  = "someQueue"
		ExpectedUrl =	"wconsj"
		hw          = test_helpers.NewWrapper()
	)
	ginkgo.BeforeEach(func() {
		if brokerDeployer != nil {
			brokerDeployer.PurgeAddressSettings()
		}
	})
	// PrepareNamespace after framework has been created
	ginkgo.JustBeforeEach(func() {
		ctx1 = sw.Framework.GetFirstContext()
		brokerDeployer = &bdw.BrokerDeploymentWrapper{}
		brokerDeployer.WithWait(true).
			WithBrokerClient(sw.BrokerClient).
			WithContext(ctx1).
			WithCustomImage(test.Config.BrokerImageName).
			WithName(DeployName).
			WithLts(!test.Config.NeedsLatestCR).
			WithConsoleExposure(true)
		hw.AddHeader("Origin", "http://localhost:8161")
        brokerDeployer.SetUpDefaultAddressSettings(AddressBit)
	})

	ginkgo.It("CollisionAvoidance check", func() {
		err := brokerDeployer.WithRedeliveryCollisionsAvoidance(AddressBit,32).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())
		urls, err := brokerDeployer.GetExternalUrls(ExpectedUrl, 0)
		address := urls[0]
		method := "org.apache.activemq.artemis:broker=\\\"0.0.0.0\\\"/getAddressSettingsAsJSON/DLQ"
		address = fmt.Sprintf("%s/console/jolokia/exec/%s", address, method)
		log.Logf("Address: %s", address)
		result, err := hw.PerformHttpRequest(address)
		gomega.Expect(err).To(gomega.BeNil())
		log.Logf("result of request: %s", result)
	})



})
