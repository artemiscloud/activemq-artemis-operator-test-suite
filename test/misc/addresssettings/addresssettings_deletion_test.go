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

	ginkgo.PIt("AutoDeleteAddresses check", func() {
		err := brokerDeployer.WithAutoDeleteAddresses(AddressBit, true).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil(), "Broker deployment failed: %s", err)

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		value := retrieveAddressSettings(address, AddressBit, hw)
		gomega.Expect(value.AutoDeleteAddresses).To(gomega.Equal(true), "AutoDeleteAddresses is not set")
	})

	ginkgo.PIt("AutoDeleteCreatedQueues check", func() {
		err := brokerDeployer.WithAutoDeleteCreatedQueues(AddressBit, true).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil(), "Broker deployment failed: %s", err)

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		value := retrieveAddressSettings(address, AddressBit, hw)
		gomega.Expect(value.AutoDeleteCreatedQueues).To(gomega.Equal(true), "AutoDeleteCreatedQueues is not set")
	})

	ginkgo.PIt("AutoDeleteQueuesMessageCount check", func() {
		err := brokerDeployer.WithAudoDeleteQueuesMessageCount(AddressBit, 100).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil(), "Broker deployment failed: %s", err)

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		value := retrieveAddressSettings(address, AddressBit, hw)
		gomega.Expect(value.AutoDeleteQueuesMessageCount).To(gomega.Equal(100), "AutoDeleteQueuesMessageCount is %d, expected 100", value.AutoDeleteQueuesMessageCount)
	})

	ginkgo.PIt("AutoDeleteJmsQueues check", func() {
		err := brokerDeployer.WithAutoDeleteJmsQueues(AddressBit, true).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil(), "Broker deployment failed: %s", err)

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		value := retrieveAddressSettings(address, AddressBit, hw)
		gomega.Expect(value.AutoDeleteJmsQueues).To(gomega.Equal(true), "AutoDeleteJmsQueues not set")
	})

	ginkgo.PIt("AutoDeleteJmsTopics check", func() {
		err := brokerDeployer.WithAutoDeleteJmsTopics(AddressBit, true).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil(), "Broker deployment failed: %s", err)

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		value := retrieveAddressSettings(address, AddressBit, hw)
		gomega.Expect(value.AutoDeleteJmsTopics).To(gomega.Equal(true), "AutoDeleteJmsTopics is not set")
	})

	ginkgo.PIt("AutoDeleteQueues check", func() {
		err := brokerDeployer.WithAutoDeleteQueues(AddressBit, true).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil(), "Broker deployment failed: %s", err)

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		value := retrieveAddressSettings(address, AddressBit, hw)
		gomega.Expect(value.AutoDeleteQueues).To(gomega.Equal(true), "AutoDeleteQueues is not set")
	})
})
