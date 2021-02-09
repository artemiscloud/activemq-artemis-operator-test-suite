package statistics

import (
	"fmt"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/framework"
	"github.com/rh-messaging/shipshape/pkg/framework/log"
	bdw "github.com/artemiscloud/activemq-artemis-operator-test-suite/pkg/bdw"
	"github.com/artemiscloud/activemq-artemis-operator-test-suite/test"
	"io/ioutil"
	"net/http"
)

var _ = ginkgo.Describe("StatisticsTest", func() {
	var (
		ctx1           *framework.ContextData
		brokerDeployer *bdw.BrokerDeploymentWrapper
	)
	//Uncomfortable bringing this to wider scope then usual.
	const (
		VarName    = "AMQ_ENABLE_METRICS_PLUGIN"
		VarValue   = "true"
		AddressBit = "metrics"
		Protocol   = test.HTTP
		// Should be available at all times
		ExpectedItem = "artemis_messages_expired"
		ExpectedUrl  = "wconsj"
	)

	//Really don't like the way its done here, but exposing this to an external wrapper isn't good either.
	testStatistics := func() error {
		gomega.Expect(brokerDeployer.DeployBrokers(1)).To(gomega.BeNil())
		brokerDeployer.SetEnvVariable(VarName, VarValue)
		log.Logf("Waiting for re-rollout of broker with updated environment")
		brokerDeployer.WaitForBrokerSet(1, 1)
		//url := test.FormUrl(Protocol, DeployName, "0", SubdomainName, ctx1.Namespace, Domain, AddressBit, Port) //nope.
		urls, err := brokerDeployer.GetExternalUrls(ExpectedUrl, 0)
		if (err!=nil) {
            return err
            
        }
		// there should be only single url in return in this case.
		url := fmt.Sprintf("%s://%s/%s/", Protocol, urls[0], AddressBit)
		resp, err := http.Get(url)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		defer resp.Body.Close()
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK))
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		bodyString := string(bodyBytes)
		// Checking for single item should be enough here.
		gomega.Expect(bodyString).To(gomega.ContainSubstring(ExpectedItem))
        return nil
	}

	// PrepareNamespace after framework has been created
	ginkgo.JustBeforeEach(func() {
		ctx1 = sw.Framework.GetFirstContext()
		brokerDeployer = &bdw.BrokerDeploymentWrapper{}
		brokerDeployer.
			WithWait(true).
			WithBrokerClient(sw.BrokerClient).
			WithContext(ctx1).
			WithCustomImage(test.Config.BrokerImageName).
			WithName(DeployName).
			WithLts(!test.Config.NeedsLatestCR)
	})

	/*
     * This tests are skipped due to ENTMQBR-3653
     * 
      ginkgo.It("StatisticsWithConsoleTestExplicit", func() {
		brokerDeployer.WithConsoleExposure(true)
		err:=testStatistics()
        gomega.Expect(err).NotTo(gomega.HaveOccurred())
        
    })

	ginkgo.It("StatisticsWithoutConsoleTestNegative", func() {
		brokerDeployer.WithConsoleExposure(false)
		err:=testStatistics()
        gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("StatisticsWithConsoleTestDefault", func() {
		brokerDeployer.WithConsoleExposure(true)
		brokerDeployer.DeployBrokers(1)
		urls, err := brokerDeployer.GetExternalUrls(ExpectedUrl, 0)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		// there should be only single url in return in this case.
		url := fmt.Sprintf("%s://%s/%s/", Protocol, urls[0], AddressBit)
        log.Logf("%s", url)
		resp, err := http.Get(url)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK))
        
	}) */
    
    ginkgo.It("StatisticsWithConsoleTestChangeSetup", func() {
        brokerDeployer.WithConsoleExposure(false)
		brokerDeployer.DeployBrokers(1)
		urls, err := brokerDeployer.GetExternalUrls(ExpectedUrl, 0)
        //No URL at this point
		gomega.Expect(err).To(gomega.HaveOccurred()) 

        brokerDeployer.WithConsoleExposure(true)
		brokerDeployer.Update()
		urls, err = brokerDeployer.GetExternalUrls(ExpectedUrl, 0)
        //Got proper URL at this point
		gomega.Expect(err).NotTo(gomega.HaveOccurred()) 

        url := fmt.Sprintf("%s://%s/%s/", Protocol, urls[0], AddressBit)
		resp, err := http.Get(url)
        bodyBytes, err := ioutil.ReadAll(resp.Body)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		bodyString := string(bodyBytes)
		// Checking for single item should be enough here.
		gomega.Expect(bodyString).To(gomega.ContainSubstring(ExpectedItem))
    })

})
