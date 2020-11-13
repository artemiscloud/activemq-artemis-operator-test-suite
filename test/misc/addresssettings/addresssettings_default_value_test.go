package addresssettings

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/framework"
	"github.com/artemiscloud/activemq-artemis-operator-test-suite/pkg/bdw"
	"github.com/artemiscloud/activemq-artemis-operator-test-suite/pkg/test_helpers"
	"github.com/artemiscloud/activemq-artemis-operator-test-suite/test"
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
		brokerDeployer.WithWait(true).
			WithBrokerClient(sw.BrokerClient).
			WithContext(ctx1).
			WithCustomImage(test.Config.BrokerImageName).
			WithName(DeployName).
			WithLts(!test.Config.NeedsLatestCR).
			WithConsoleExposure(true)
		brokerDeployer.SetUpDefaultAddressSettings(AddressBit)
	})

	ginkgo.It("DefaultConsumer check", func() {
		err := brokerDeployer.WithDefaultConsumerBeforeDispatch(AddressBit, 1).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		verifyAddressSettingsInt(address, AddressBit, "defaultConsumersBeforeDispatch", 1, hw)
	})

	ginkgo.It("DefaultConsumerWindowSize check", func() {
		err := brokerDeployer.WithDefaultConsumerWindowSize(AddressBit, 1234567).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		verifyAddressSettingsInt(address, AddressBit, "defaultConsumerWindowSize", 1234567, hw)
	})

	ginkgo.It("DelayBeforeDispatch check", func() {
		err := brokerDeployer.WithDefaultDelayBeforeDispatch(AddressBit, 150).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		verifyAddressSettingsInt(address, AddressBit, "defaultDelayBeforeDispatch", 150, hw)
	})

	ginkgo.It("DefaultExclusiveQueue check", func() {
		err := brokerDeployer.WithDefaultExclusiveQueue(AddressBit, true).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		verifyAddressSettingsBool(address, AddressBit, "defaultExclusiveQueue", true, hw)
	})

	ginkgo.It("DefaultGroupBuckets check", func() {
		err := brokerDeployer.WithDefaultGroupBuckets(AddressBit, 10).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		verifyAddressSettingsInt(address, AddressBit, "defaultGroupBuckets", 10, hw)
	})

	ginkgo.It("DefaultGroupFirstKey check", func() {
		err := brokerDeployer.WithDefaultGroupFirstKey(AddressBit, "hey").DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		verifyAddressSettingsString(address, AddressBit, "defaultGroupFirstKey", "hey", hw)
	})

	ginkgo.It("DefaultGroupRebalance check", func() {
		err := brokerDeployer.WithDefaultGroupRebalance(AddressBit, true).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		verifyAddressSettingsBool(address, AddressBit, "defaultGroupRebalance", true, hw)
	})

	/* // TODO: This is NOT expected to work due to issue in the init container
	    ginkgo.It("DefaultGroupRebalancePauseDispatch check", func() {
			err := brokerDeployer.WithDefaultGroupRebalancePauseDispatch(AddressBit, true).DeployBrokers(1)
			gomega.Expect(err).To(gomega.BeNil())

			urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
			address := urls[0]
			domain := strings.Split(address, ".")[0]
			header := strings.Replace(OriginHeader, "NAME", domain, 1)
			hw.AddHeader("Origin", header)
			actualUrl := "http://admin:admin@" + address + CallAddress + AddressBit
			hw.WithPassword("admin").WithUser("admin")
			result, err := hw.PerformHttpRequest(actualUrl)
			if err != nil {
				log.Logf("%s", err)
			}
			var item map[string]map[string]string
			json.Unmarshal([]byte(result), &item)
	        brokerValue := item["value"]["defaultGroupRebalancePauseDispatch"]
			gomega.Expect(strconv.ParseBool(brokerValue)).To(gomega.Equal(true))
		})
	*/

	ginkgo.It("DefaultLastValueKey check", func() {
		err := brokerDeployer.WithDefaultLastValueKey(AddressBit, "hey").DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		verifyAddressSettingsString(address, AddressBit, "defaultLastValueKey", "hey", hw)
	})

	ginkgo.It("DefaultLastValueQueue check", func() {
		err := brokerDeployer.WithDefaultLastValueQueue(AddressBit, true).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		verifyAddressSettingsBool(address, AddressBit, "defaultLastValueQueue", true, hw)
	})

	ginkgo.It("DefaultMaxConsumers check", func() {
		err := brokerDeployer.WithDefaultMaxConsumers(AddressBit, 32).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		urls, err := brokerDeployer.GetExternalUrls(ExpectedURL, 0)
		address := urls[0]
		verifyAddressSettingsInt(address, AddressBit, "defaultMaxConsumers", 32, hw)
	})
})
