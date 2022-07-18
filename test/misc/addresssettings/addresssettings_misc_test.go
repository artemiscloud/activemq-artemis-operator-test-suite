package addresssettings

import (
	"github.com/artemiscloud/activemq-artemis-operator-test-suite/pkg/bdw"
	"github.com/artemiscloud/activemq-artemis-operator-test-suite/pkg/test_helpers"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/framework"
)

var _ = ginkgo.Describe("AddressSettingMiscTest", func() {

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

	ginkgo.PIt("DLQPrefix check", func() {
		err := brokerDeployer.WithDlqPrefix(AddressBit, "prefix").DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil(), "Broker deployment failed: %s", err)

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		value := retrieveAddressSettings(address, AddressBit, hw)
		gomega.Expect(value.DeadLetterQueuePrefix).To(gomega.Equal("prefix"), "DeadLetterQueuePrefix is \"%s\", expected \"prefix\"", value.DeadLetterQueuePrefix)
	})

	ginkgo.PIt("DLQSuffix check", func() {
		err := brokerDeployer.WithDlqSuffix(AddressBit, "suffix").DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil(), "Broker deployment failed: %s", err)

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		value := retrieveAddressSettings(address, AddressBit, hw)
		gomega.Expect(value.DeadLetterQueueSuffix).To(gomega.Equal("suffix"), "DeadLetterQueueSuffix is \"%s\", expected \"suffix\"", value.DeadLetterQueueSuffix)
	})

	ginkgo.PIt("DLQAddress check", func() {
		err := brokerDeployer.WithDeadLetterAddress(AddressBit, "DLqQ").DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil(), "Broker deployment failed: %s", err)

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		value := retrieveAddressSettings(address, AddressBit, hw)
		gomega.Expect(value.DLA).To(gomega.Equal("DLqQ"), "DLA is \"%s\", expected \"DLqQ\"", value.DLA)
	})

	ginkgo.PIt("AddressFullPolicy check", func() {
		err := brokerDeployer.WithAddressFullPolicy(AddressBit, bdw.DropPolicy).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil(), "Broker deployment failed: %s", err)

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		value := retrieveAddressSettings(address, AddressBit, hw)
		gomega.Expect(value.AddressFullMessagePolicy).To(gomega.Equal("DROP"), "AddressFullMessagePolicy is \"%s\", expected \"DROP\"", value.AddressFullMessagePolicy)
	})

	ginkgo.PIt("MetricsCheck check", func() {
		err := brokerDeployer.WithEnableMetrics(AddressBit, true).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil(), "Broker deployment failed: %s", err)

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		value := retrieveAddressSettings(address, AddressBit, hw)
		gomega.Expect(value.EnableMetrics).To(gomega.Equal(true), "EnableMetrics is not set")
	})

	/*
		ginkgo.PIt("MetricsCheck check", func() {
			err := brokerDeployer.WithManagementBrowsePageSize(AddressBit, 101).DeployBrokers(1)
			gomega.Expect(err).To(gomega.BeNil())

			urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
			address := urls[0]
			verifyAddressSettingsBool(address, AddressBit, "enableMetrics",true, hw)
		})
	*/

	ginkgo.PIt("SlowConsumerCheck check", func() {
		err := brokerDeployer.WithSlowConsumerCheckPeriod(AddressBit, 10).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil(), "Broker deployment failed: %s", err)

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		gomega.Expect(err).To(gomega.BeNil(), "Can not receive URLs from openshift: %s", err)
		address := urls[0]
		value := retrieveAddressSettings(address, AddressBit, hw)
		gomega.Expect(value.SlowConsumerCheckPeriod).To(gomega.Equal(10), "SlowConsumerCheckPeriod is %d, expected 10", value.SlowConsumerCheckPeriod)
	})

	ginkgo.PIt("SlowConsumerPolicy check", func() {
		err := brokerDeployer.WithSlowConsumerPolicy(AddressBit, bdw.Notify).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil(), "Broker deployment failed: %s", err)

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		gomega.Expect(err).To(gomega.BeNil(), "Can not receive URLs from openshift: %s", err)
		address := urls[0]
		value := retrieveAddressSettings(address, AddressBit, hw)
		gomega.Expect(value.SlowConsumerPolicy).To(gomega.Equal(bdw.NOTIFY), "SlowConsumerPolicy is %s, expected NOTIFY", value.SlowConsumerPolicy)
	})

	ginkgo.PIt("SlowConsumerThreshold check", func() {
		err := brokerDeployer.WithSlowConsumerThreshold(AddressBit, 320).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil(), "Broker deployment failed: %s", err)

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		gomega.Expect(err).To(gomega.BeNil(), "Can not receive URLs from openshift: %s", err)
		address := urls[0]
		value := retrieveAddressSettings(address, AddressBit, hw)
		gomega.Expect(value.SlowConsumerThreshold).To(gomega.Equal(320), "SlowConsumerThreshold is %d, expected 320", value.SlowConsumerThreshold)
	})
})
