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

	ginkgo.It("UpdateAddressSetting check", func() {
		err := brokerDeployer.WithAutoCreateAddresses(AddressBit, true).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		value := retrieveAddressSettings(address, AddressBit, hw)
		gomega.Expect(value.AutoCreateAddresses).To(gomega.Equal(true), "AddressAutoCreation setting is not correct")

		brokerDeployer.WithAutoCreateAddresses(AddressBit, false)
		err = brokerDeployer.Update()
		gomega.Expect(err).To(gomega.BeNil(), "error during Broker addressettings update on openshift")

		value = retrieveAddressSettings(address, AddressBit, hw)
		gomega.Expect(value.AutoCreateAddresses).To(gomega.Equal(false), "Addressettings item was not updated")
	})

})
