package addresssettings

import (
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

	ginkgo.It("CollisionAvoidance check", func() {
		err := brokerDeployer.WithRedeliveryCollisionsAvoidance(AddressBit, 1).DeployBrokers(1)
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
        
		brokerValue := item["value"]["redeliveryCollisionAvoidanceFactor"]
		gomega.Expect(brokerValue).To(gomega.Equal(string(1.0f)))
		
	})
    
   
    ginkgo.It("RedeliveryDelayMultiplier check", func() {
		err := brokerDeployer.WithRedeliveryDelayMult(AddressBit,1).DeployBrokers(1)
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
        brokerValue := item["value"]["redeliveryMultiplier"]
		gomega.Expect(brokerValue).To(gomega.Equal(string(1.0f)))
	})

	ginkgo.It("RedeliveryDelay check", func() {
		err := brokerDeployer.WithRedeliveryDelay(AddressBit, 1).DeployBrokers(1)
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
        brokerValue := item["value"]["redeliveryDelay"]
		gomega.Expect(brokerValue).To(gomega.Equal(string(1)))
	})

})
