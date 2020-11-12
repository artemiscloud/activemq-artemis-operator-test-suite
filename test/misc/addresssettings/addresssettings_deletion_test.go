package addresssettings

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/framework"
	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/pkg/bdw"
	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/pkg/test_helpers"
	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/test"
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

	ginkgo.It("AutoDeleteAddresses check", func() {
		err := brokerDeployer.WithAudoDeleteAddresses(AddressBit, true).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		verifyAddressSettingsBool(address, AddressBit, "autoDeleteAddresses", true, hw)
	})

	ginkgo.It("AutoDeleteCreatedQueues check", func() {
		err := brokerDeployer.WithAutoDeleteCreatedQueues(AddressBit, true).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		verifyAddressSettingsBool(address, AddressBit, "autoDeleteCreatedQueues", true, hw)
	})

	ginkgo.It("AutoDeleteQueuesMessageCount check", func() {
		err := brokerDeployer.WithAudoDeleteQueuesMessageCount(AddressBit, 100).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		verifyAddressSettingsInt(address, AddressBit, "autoDeleteQueuesMessageCount", 100, hw)
	})

	ginkgo.It("AutoDeleteJmsQueues check", func() {
		err := brokerDeployer.WithAutoDeleteJmsQueues(AddressBit, true).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		verifyAddressSettingsBool(address, AddressBit, "autoDeleteJmsQueues", true, hw)
	})

	ginkgo.It("AutoDeleteJmsTopics check", func() {
		err := brokerDeployer.WithAutoDeleteJmsTopics(AddressBit, true).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		verifyAddressSettingsBool(address, AddressBit, "autoDeleteJmsTopics", true, hw)
	})

	ginkgo.It("AutoDeleteQueues check", func() {
		err := brokerDeployer.WithAutoDeleteQueues(AddressBit, true).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		verifyAddressSettingsBool(address, AddressBit, "autoDeleteQueues", true, hw)
	})
})
