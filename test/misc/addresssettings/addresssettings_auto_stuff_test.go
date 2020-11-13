package addresssettings

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/framework"
	"github.com/artemiscloud/activemq-artemis-operator-test-suite/pkg/bdw"
	"github.com/artemiscloud/activemq-artemis-operator-test-suite/pkg/test_helpers"
	"github.com/artemiscloud/activemq-artemis-operator-test-suite/test"
)

var _ = ginkgo.Describe("AddressSettingsDeletionTest", func() {

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
		brokerDeployer.WithWait(true).
			WithBrokerClient(sw.BrokerClient).
			WithContext(ctx1).
			WithCustomImage(test.Config.BrokerImageName).
			WithName(DeployName).
			WithLts(!test.Config.NeedsLatestCR).
			WithConsoleExposure(true)
		brokerDeployer.SetUpDefaultAddressSettings(AddressBit)
	})

	ginkgo.It("AutoCreateAddresses check", func() {
		err := brokerDeployer.WithAutoCreateAddresses(AddressBit, true).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		verifyAddressSettingsBool(address, AddressBit, "autoCreateAddresses", true, hw)
	})

	ginkgo.It("AutoCreateDeadLetterResources check", func() {
		err := brokerDeployer.WithAutoCreateDeadLetterResources(AddressBit, true).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		verifyAddressSettingsBool(address, AddressBit, "autoCreateDeadLetterResources", true, hw)
	})

	ginkgo.It("AutoCreateExpiryResources check", func() {
		err := brokerDeployer.WithAutoCreateExpiryResources(AddressBit, true).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		verifyAddressSettingsBool(address, AddressBit, "autoCreateExpiryResources", true, hw)
	})

	ginkgo.It("AutoCreateJmsQueues check", func() {
		err := brokerDeployer.WithAutoCreateJmsQueues(AddressBit, true).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		verifyAddressSettingsBool(address, AddressBit, "autoCreateJmsQueues", true, hw)
	})

	ginkgo.It("AutoCreateJmsTopics check", func() {
		err := brokerDeployer.WithAutoCreateJmsTopics(AddressBit, true).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		verifyAddressSettingsBool(address, AddressBit, "autoCreateJmsTopics", true, hw)
	})

	ginkgo.It("AutoCreateQueues check", func() {
		err := brokerDeployer.WithAutoCreateQueues(AddressBit, true).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		verifyAddressSettingsBool(address, AddressBit, "autoCreateQueues", true, hw)
	})
})
