package addresssettings

import (
	"github.com/artemiscloud/activemq-artemis-operator-test-suite/pkg/bdw"
	"github.com/artemiscloud/activemq-artemis-operator-test-suite/pkg/test_helpers"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/framework"
)

var _ = ginkgo.Describe("AddressSettingsDefaultValueTest", func() {

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

	ginkgo.It("DefaultConsumer check", func() {
		err := brokerDeployer.WithDefaultConsumerBeforeDispatch(AddressBit, 1).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil(), "Broker deployment failed: %s", err)

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		value := retrieveAddressSettings(address, AddressBit, hw)
		gomega.Expect(value.DefaultConsumersBeforeDispatch).To(gomega.Equal(1), "DefaultConsumersBeforeDispatch is %d, expected 1", value.DefaultConsumersBeforeDispatch)
	})

	ginkgo.It("DefaultConsumerWindowSize check", func() {
		err := brokerDeployer.WithDefaultConsumerWindowSize(AddressBit, 1234567).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil(), "Broker deployment failed: %s", err)

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		value := retrieveAddressSettings(address, AddressBit, hw)
		gomega.Expect(value.DefaultConsumerWindowSize).To(gomega.Equal(1234567), "DefaultConsumerWindowSize is %d, expected 1234567", value.DefaultConsumerWindowSize)
	})

	ginkgo.It("DelayBeforeDispatch check", func() {
		err := brokerDeployer.WithDefaultDelayBeforeDispatch(AddressBit, 150).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil(), "Broker deployment failed: %s", err)

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		value := retrieveAddressSettings(address, AddressBit, hw)
		gomega.Expect(value.DefaultDelayBeforeDispatch).To(gomega.Equal(150), "DefaultDelayBeforeDispatch is %s, expected 150", value.DefaultDelayBeforeDispatch)
	})

	ginkgo.It("DefaultExclusiveQueue check", func() {
		err := brokerDeployer.WithDefaultExclusiveQueue(AddressBit, true).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil(), "Broker deployment failed: %s", err)

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		value := retrieveAddressSettings(address, AddressBit, hw)
		gomega.Expect(value.DefaultExclusiveQueue).To(gomega.Equal(true), "DefaultExclusive queue not set")
	})

	ginkgo.It("DefaultGroupBuckets check", func() {
		err := brokerDeployer.WithDefaultGroupBuckets(AddressBit, 10).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil(), "Broker deployment failed: %s", err)

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		value := retrieveAddressSettings(address, AddressBit, hw)
		gomega.Expect(value.DefaultGroupBuckets).To(gomega.Equal(10), "DefaultGroupBuckets is, expected 10", value.DefaultGroupBuckets)
	})

	ginkgo.It("DefaultGroupFirstKey check", func() {
		err := brokerDeployer.WithDefaultGroupFirstKey(AddressBit, "hey").DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		value := retrieveAddressSettings(address, AddressBit, hw)
		gomega.Expect(value.DefaultGroupFirstKey).To(gomega.Equal("hey"), "DefaultGroupFirstKey is \"%s\", expected \"hey\"", value.DefaultGroupFirstKey)
	})

	ginkgo.It("DefaultGroupRebalance check", func() {
		err := brokerDeployer.WithDefaultGroupRebalance(AddressBit, true).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil(), "Broker deployment failed: %s", err)

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		value := retrieveAddressSettings(address, AddressBit, hw)
		gomega.Expect(value.DefaultGroupRebalance).To(gomega.Equal(true), "DefaultGroupRebalance not set")
	})

	ginkgo.It("DefaultLastValueKey check", func() {
		err := brokerDeployer.WithDefaultLastValueKey(AddressBit, "hey").DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil(), "Broker deployment failed: %s", err)

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		value := retrieveAddressSettings(address, AddressBit, hw)
		gomega.Expect(value.DefaultLastValueKey).To(gomega.Equal("hey"), "DefaultLastValueKey is \"%s\", expected \"hey\"", value.DefaultLastValueKey)
	})

	ginkgo.It("DefaultLastValueQueue check", func() {
		err := brokerDeployer.WithDefaultLastValueQueue(AddressBit, true).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil(), "Broker deployment failed: %s", err)

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		value := retrieveAddressSettings(address, AddressBit, hw)
		gomega.Expect(value.LastValueQueue).To(gomega.Equal(true), "LastValueQueue is not set")
	})

	ginkgo.It("DefaultMaxConsumers check", func() {
		err := brokerDeployer.WithDefaultMaxConsumers(AddressBit, 32).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil(), "Broker deployment failed: %s", err)

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		value := retrieveAddressSettings(address, AddressBit, hw)
		gomega.Expect(value.DefaultMaxConsumers).To(gomega.Equal(32), "DefaultMaxConsumers is %d, expected 32", value.DefaultMaxConsumers)
	})
})
