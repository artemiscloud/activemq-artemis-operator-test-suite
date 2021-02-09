package addresssettings

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/framework"
	"github.com/artemiscloud/activemq-artemis-operator-test-suite/pkg/bdw"
	"github.com/artemiscloud/activemq-artemis-operator-test-suite/pkg/test_helpers"
	"github.com/artemiscloud/activemq-artemis-operator-test-suite/test"
)

var _ = ginkgo.Describe("AddressSettingsExpiryCheck", func() {

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

	ginkgo.It("ExpiryAddress check", func() {
		err := brokerDeployer.WithExpiryAddress(AddressBit, "expire").DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		value:= retrieveAddressSettings(address,AddressBit, hw)
        gomega.Expect(value.ExpiryAddress).To(gomega.Equal("expire"))
	})

	ginkgo.It("ExpiryDelay check", func() {
		err := brokerDeployer.WithExpiryDelay(AddressBit, 1).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		value:= retrieveAddressSettings(address,AddressBit, hw)
        gomega.Expect(value.ExpiryDelay).To(gomega.Equal(1))
	})

	ginkgo.It("ExpiryPrefix check", func() {
		err := brokerDeployer.WithExpiryPrefix(AddressBit, "prefix").DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		value:= retrieveAddressSettings(address,AddressBit, hw)
        gomega.Expect(value.ExpiryQueuePrefix).To(gomega.Equal("prefix"))
	})

	ginkgo.It("ExpirySuffix check", func() {
		err := brokerDeployer.WithExpirySuffix(AddressBit, "suffix").DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		value:= retrieveAddressSettings(address,AddressBit, hw)
        gomega.Expect(value.ExpiryQueueSuffix).To(gomega.Equal("suffix"))
	})

	ginkgo.It("MinExpiryDelay check", func() {
		err := brokerDeployer.WithMinExpiryDelay(AddressBit, 101).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		value:= retrieveAddressSettings(address,AddressBit, hw)
        gomega.Expect(value.MinExpiryDelay).To(gomega.Equal(101))
	})

	ginkgo.It("MaxExpiryDelay check", func() {
		err := brokerDeployer.WithMaxExpiryDelay(AddressBit, 101).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		value:= retrieveAddressSettings(address,AddressBit, hw)
        gomega.Expect(value.MaxExpiryDelay).To(gomega.Equal(101))
	})
})
