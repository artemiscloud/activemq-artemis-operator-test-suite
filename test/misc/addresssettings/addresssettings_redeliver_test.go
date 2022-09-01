package addresssettings

import (
	"github.com/artemiscloud/activemq-artemis-operator-test-suite/pkg/bdw"
	"github.com/artemiscloud/activemq-artemis-operator-test-suite/pkg/test_helpers"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/framework"
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
		ExpectedURL = "wconsj"
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
		setEnv(ctx1, brokerDeployer)
		brokerDeployer.SetUpDefaultAddressSettings(AddressBit)
	})

	ginkgo.PIt("CollisionAvoidance check", func() {
		err := brokerDeployer.WithRedeliveryCollisionsAvoidance(AddressBit, "1").DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil(), "Broker deployment failed: %s", err)
		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		gomega.Expect(err).To(gomega.BeNil(), "Can not retrieve URLs from openshift: %s", err)
		address := urls[0]
		value := retrieveAddressSettings(address, AddressBit, hw)
		gomega.Expect(value.RedeliveryCollisionAvoidanceFactor).To(gomega.Equal(float32(1.0)), "RedeliveryCollisionAvoidanceFactor is %f, expected 1.0", value.RedeliveryCollisionAvoidanceFactor)
	})

	ginkgo.PIt("RedeliveryDelayMultiplier check", func() {
		err := brokerDeployer.WithRedeliveryDelayMult(AddressBit, "1").DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil(), "Broker deployment failed: %s", err)

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		gomega.Expect(err).To(gomega.BeNil(), "Can not retrieve URLs from openshift: %s", err)
		address := urls[0]
		value := retrieveAddressSettings(address, AddressBit, hw)
		gomega.Expect(value.RedeliveryMultiplier).To(gomega.Equal(float32(1.0)), "RedeliveryMultiplier is %f, expected 1.0", value.RedeliveryMultiplier)
	})

	ginkgo.PIt("RedeliveryDelay check", func() {
		err := brokerDeployer.WithRedeliveryDelay(AddressBit, 1).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil(), "Broker deployment failed: %s", err)

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		gomega.Expect(err).To(gomega.BeNil(), "Can not retrieve URLs from openshift: %s", err)
		address := urls[0]
		value := retrieveAddressSettings(address, AddressBit, hw)
		gomega.Expect(value.RedeliveryDelay).To(gomega.Equal(1), "RedeliveryDelay is %d, expected 1", value.RedeliveryDelay)
	})
})
