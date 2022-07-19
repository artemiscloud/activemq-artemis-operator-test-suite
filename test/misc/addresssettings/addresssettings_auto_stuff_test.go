package addresssettings

import (
	"github.com/artemiscloud/activemq-artemis-operator-test-suite/pkg/bdw"
	"github.com/artemiscloud/activemq-artemis-operator-test-suite/pkg/test_helpers"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/framework"
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
		setEnv(ctx1, brokerDeployer)
		brokerDeployer.SetUpDefaultAddressSettings(AddressBit)
	})

	ginkgo.PIt("AutoCreateAddresses check", func() {
		err := brokerDeployer.WithAutoCreateAddresses(AddressBit, true).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil(), "Broker deployment failed: %s", err)

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		value := retrieveAddressSettings(address, AddressBit, hw)
		gomega.Expect(value.AutoCreateAddresses).To(gomega.Equal(true), "AutoCreateAddresses not set")
	})

	ginkgo.PIt("AutoCreateDeadLetterResources check", func() {
		err := brokerDeployer.WithAutoCreateDeadLetterResources(AddressBit, true).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil(), "Broker deployment failed: %s", err)

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		value := retrieveAddressSettings(address, AddressBit, hw)
		gomega.Expect(value.AutoCreateDeadLetterResources).To(gomega.Equal(true), "AutoCreateDeadLetterResources not set")
	})

	ginkgo.PIt("AutoCreateExpiryResources check", func() {
		err := brokerDeployer.WithAutoCreateExpiryResources(AddressBit, true).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil(), "Broker deployment failed: %s", err)

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		value := retrieveAddressSettings(address, AddressBit, hw)
		gomega.Expect(value.AutoCreateExpiryResources).To(gomega.Equal(true), "AutoCreateExpiry resources not set")
	})

	ginkgo.PIt("AutoCreateJmsQueues check", func() {
		err := brokerDeployer.WithAutoCreateJmsQueues(AddressBit, true).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil(), "Broker deployment failed: %s", err)

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		value := retrieveAddressSettings(address, AddressBit, hw)
		gomega.Expect(value.AutoCreateJmsQueues).To(gomega.Equal(true), "AutoCreateJmsQueues not set")
	})

	ginkgo.PIt("AutoCreateJmsTopics check", func() {
		err := brokerDeployer.WithAutoCreateJmsTopics(AddressBit, true).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil(), "Broker deployment failed: %s", err)

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		value := retrieveAddressSettings(address, AddressBit, hw)
		gomega.Expect(value.AutoCreateJmsTopics).To(gomega.Equal(true), "AutoCreateJmstopics not set")
	})

	ginkgo.PIt("AutoCreateQueues check", func() {
		err := brokerDeployer.WithAutoCreateQueues(AddressBit, true).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil(), "Broker deployment failed: %s", err)

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		value := retrieveAddressSettings(address, AddressBit, hw)
		gomega.Expect(value.AutoCreateQueues).To(gomega.Equal(true), "AutoCreateQueues not set")
	})
})
