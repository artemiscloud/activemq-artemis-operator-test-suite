package addresssettings

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/framework"
	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/pkg/bdw"
	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/pkg/test_helpers"
	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/test"
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
		ExpectedUrl = "wconsj"
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

		urls, err := brokerDeployer.GetExternalUrls(ExpectedUrl, 0)
		address := urls[0]
        verifyAddressSettingsString(address, AddressBit, "expiryAddress","expire", hw)
	})
    
   
    ginkgo.It("ExpiryDelay check", func() {
		err := brokerDeployer.WithExpiryDelay(AddressBit,1).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())
        
		urls, err := brokerDeployer.GetExternalUrls(ExpectedUrl, 0)
		address := urls[0]
        verifyAddressSettingsInt(address, AddressBit, "expiryDelay",1, hw)
	})

	ginkgo.It("ExpiryPrefix check", func() {
		err := brokerDeployer.WithExpiryPrefix(AddressBit, "prefix").DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		urls, err := brokerDeployer.GetExternalUrls(ExpectedUrl, 0)
		address := urls[0]
        verifyAddressSettingsString(address, AddressBit, "expiryQueuePrefix","prefix", hw)
	})
        
    ginkgo.It("ExpirySuffix check", func() {
		err := brokerDeployer.WithExpirySuffix(AddressBit, "suffix").DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		urls, err := brokerDeployer.GetExternalUrls(ExpectedUrl, 0)
		address := urls[0]
        verifyAddressSettingsString(address, AddressBit, "expiryQueueSuffix","suffix", hw)
	})
    
    ginkgo.It("MinExpiryDelay check", func() {
		err := brokerDeployer.WithMinExpiryDelay(AddressBit, 101).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		urls, err := brokerDeployer.GetExternalUrls(ExpectedUrl, 0)
		address := urls[0]
        verifyAddressSettingsInt(address, AddressBit, "minExpiryDelay",101, hw)
	})
    ginkgo.It("MaxExpiryDelay check", func() {
		err := brokerDeployer.WithMaxExpiryDelay(AddressBit, 101).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		urls, err := brokerDeployer.GetExternalUrls(ExpectedUrl, 0)
		address := urls[0]
        verifyAddressSettingsInt(address, AddressBit, "maxExpiryDelay",101, hw)
	})
})
