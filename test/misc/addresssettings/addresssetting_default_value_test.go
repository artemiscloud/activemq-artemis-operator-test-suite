package addresssettings

import (
    "strconv"
	"encoding/json"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/framework"
	"github.com/rh-messaging/shipshape/pkg/framework/log"
	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/pkg/bdw"
	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/pkg/test_helpers"
	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/test"
	"strings"
)

var _ = ginkgo.Describe("AddressSettingsRedeliveryTest", func() {

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

	ginkgo.It("DefaultConsumer check", func() {
		err := brokerDeployer.WithDefaultConsumerBeforeDispatch(AddressBit, 1).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		urls, err := brokerDeployer.GetExternalUrls(ExpectedUrl, 0)
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
        
		brokerValue := item["value"]["defaultConsumersBeforeDispatch"]
		gomega.Expect(brokerValue).To(gomega.Equal(string(1)))
	})
    
   
    ginkgo.It("DefaultConsumerWindowSize check", func() {
		err := brokerDeployer.WithDefaultConsumerWindowSize(AddressBit,1234567).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())
        
		urls, err := brokerDeployer.GetExternalUrls(ExpectedUrl, 0)
		address := urls[0]
        domain := strings.Split(address, ".")[0]
        header := strings.Replace(OriginHeader,"NAME", domain,1)
        hw.AddHeader("Origin", header)
        actualUrl := "http://admin:admin@"+ address + CallAddress + AddressBit
        hw.WithPassword("admin").WithUser("admin")
        result, err := hw.PerformHttpRequest(actualUrl)
        if err != nil {
            log.Logf("%s", err) 
        }
        var item map[string]map[string]string
        json.Unmarshal([]byte(result), &item)
        brokerValue := item["value"]["defaultConsumerWindowSize"]
		gomega.Expect(brokerValue).To(gomega.Equal(string(1234567)))
	})

	ginkgo.It("DelayBeforeDispatch check", func() {
		err := brokerDeployer.WithDefaultDelayBeforeDispatch(AddressBit, 150).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		urls, err := brokerDeployer.GetExternalUrls(ExpectedUrl, 0)
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
        brokerValue := item["value"]["defaultDelayBeforeDispatch"]
		gomega.Expect(brokerValue).To(gomega.Equal(string(150)))
	})
    
    ginkgo.It("DefaultExclusiveQueue check", func() {
		err := brokerDeployer.WithDefaultExclusiveQueue(AddressBit, true).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		urls, err := brokerDeployer.GetExternalUrls(ExpectedUrl, 0)
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
        brokerValue := item["value"]["defaultExclusiveQueue"]
		gomega.Expect(strconv.ParseBool(brokerValue)).To(gomega.Equal(true))
	})
    
    ginkgo.It("DefaultGroupBuckets check", func() {
		err := brokerDeployer.WithDefaultGroupBuckets(AddressBit, 10).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		urls, err := brokerDeployer.GetExternalUrls(ExpectedUrl, 0)
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
        brokerValue := item["value"]["defaultGroupBuckets"]
		gomega.Expect(brokerValue).To(gomega.Equal(string(10)))
	})
    
    ginkgo.It("DefaultGroupFirstKey check", func() {
		err := brokerDeployer.WithDefaultGroupFirstKey(AddressBit, "hey").DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		urls, err := brokerDeployer.GetExternalUrls(ExpectedUrl, 0)
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
        brokerValue := item["value"]["defaultGroupFirstKey"]
		gomega.Expect(brokerValue).To(gomega.Equal("hey"))
	})
        
    ginkgo.It("DefaultGroupRebalance check", func() {
		err := brokerDeployer.WithDefaultGroupRebalance(AddressBit, true).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		urls, err := brokerDeployer.GetExternalUrls(ExpectedUrl, 0)
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
        brokerValue := item["value"]["defaultGroupRebalance"]
		gomega.Expect(strconv.ParseBool(brokerValue)).To(gomega.Equal(true))
	})
   /* // TODO: This is NOT expected to work due to issue in the init container
    ginkgo.It("DefaultGroupRebalancePauseDispatch check", func() {
		err := brokerDeployer.WithDefaultGroupRebalancePauseDispatch(AddressBit, true).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		urls, err := brokerDeployer.GetExternalUrls(ExpectedUrl, 0)
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
	}) */

    ginkgo.It("DefaultLastValueKey check", func() {
		err := brokerDeployer.WithDefaultLastValueKey(AddressBit, "hey").DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		urls, err := brokerDeployer.GetExternalUrls(ExpectedUrl, 0)
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
        brokerValue := item["value"]["defaultLastValueKey"]
		gomega.Expect(brokerValue).To(gomega.Equal("hey"))
	})
    
    ginkgo.It("DefaultLastValueQueue check", func() {
		err := brokerDeployer.WithDefaultLastValueQueue(AddressBit, true).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		urls, err := brokerDeployer.GetExternalUrls(ExpectedUrl, 0)
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
        brokerValue := item["value"]["defaultLastValueQueue"]
		gomega.Expect(strconv.ParseBool(brokerValue)).To(gomega.Equal(true))
	})
    
    ginkgo.It("DefaultMaxConsumers check", func() {
		err := brokerDeployer.WithDefaultMaxConsumers(AddressBit, 32).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		urls, err := brokerDeployer.GetExternalUrls(ExpectedUrl, 0)
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
        brokerValue := item["value"]["defaultMaxConsumers"]
		gomega.Expect(brokerValue).To(gomega.Equal(32))
	})

})
