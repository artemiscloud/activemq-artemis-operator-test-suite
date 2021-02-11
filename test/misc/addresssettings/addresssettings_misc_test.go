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

	ginkgo.It("DLQPrefix check", func() {
		err := brokerDeployer.WithDlqPrefix(AddressBit, "prefix").DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		value := retrieveAddressSettings(address, AddressBit, hw)
		gomega.Expect(value.DeadLetterQueuePrefix).To(gomega.Equal("prefix"))
	})

	ginkgo.It("DLQSuffix check", func() {
		err := brokerDeployer.WithDlqSuffix(AddressBit, "suffix").DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		value := retrieveAddressSettings(address, AddressBit, hw)
		gomega.Expect(value.DeadLetterQueueSuffix).To(gomega.Equal("suffix"))
	})

	ginkgo.It("DLQAddress check", func() {
		err := brokerDeployer.WithDeadLetterAddress(AddressBit, "DLqQ").DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		value := retrieveAddressSettings(address, AddressBit, hw)
		gomega.Expect(value.DLA).To(gomega.Equal("DLqQ"))
	})

	ginkgo.It("AddressFullPolicy check", func() {
		err := brokerDeployer.WithAddressFullPolicy(AddressBit, bdw.DropPolicy).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		value := retrieveAddressSettings(address, AddressBit, hw)
		gomega.Expect(value.AddressFullMessagePolicy).To(gomega.Equal("DROP"))
	})

	ginkgo.It("MetricsCheck check", func() {
		err := brokerDeployer.WithEnableMetrics(AddressBit, true).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		value := retrieveAddressSettings(address, AddressBit, hw)
		gomega.Expect(value.EnableMetrics).To(gomega.Equal(true))
	})

	/*
		ginkgo.It("MetricsCheck check", func() {
			err := brokerDeployer.WithManagementBrowsePageSize(AddressBit, 101).DeployBrokers(1)
			gomega.Expect(err).To(gomega.BeNil())

			urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
			address := urls[0]
			verifyAddressSettingsBool(address, AddressBit, "enableMetrics",true, hw)
		})
	*/

	ginkgo.It("SlowConsumerCheck check", func() {
		err := brokerDeployer.WithSlowConsumerCheckPeriod(AddressBit, 10).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		value := retrieveAddressSettings(address, AddressBit, hw)
		gomega.Expect(value.SlowConsumerCheckPeriod).To(gomega.Equal(10))
	})

	ginkgo.It("SlowConsumerPolicy check", func() {
		err := brokerDeployer.WithSlowConsumerPolicy(AddressBit, bdw.Notify).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		value := retrieveAddressSettings(address, AddressBit, hw)
		gomega.Expect(value.SlowConsumerPolicy).To(gomega.Equal(bdw.NOTIFY))
	})

	ginkgo.It("SlowConsumerThreshold check", func() {
		err := brokerDeployer.WithSlowConsumerThreshold(AddressBit, 320).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		value := retrieveAddressSettings(address, AddressBit, hw)
		gomega.Expect(value.SlowConsumerThreshold).To(gomega.Equal(320))
	})
})
