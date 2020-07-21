package configuration

import (
	"fmt"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/framework"
	"github.com/rh-messaging/shipshape/pkg/framework/log"
	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/test"
	"io/ioutil"
	"net/http"
)

var _ = ginkgo.Describe("MetricsTest", func() {
	var (
		ctx1 *framework.ContextData
		bdw  *test.BrokerDeploymentWrapper
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
	testStatistics := func() {
		gomega.Expect(bdw.DeployBrokers(1)).To(gomega.BeNil())
		bdw.SetEnvVariable(VarName, VarValue)
		log.Logf("Waiting for re-rollout of broker with updated environment")
		bdw.WaitForBrokerSet(1, 1)
		//url := test.FormUrl(Protocol, DeployName, "0", SubdomainName, ctx1.Namespace, Domain, AddressBit, Port) //nope.
		urls, err := bdw.GetExternalUrls(ExpectedUrl, 0)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
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
	}

	// PrepareNamespace after framework has been created
	ginkgo.JustBeforeEach(func() {
		ctx1 = sw.Framework.GetFirstContext()
		bdw = &test.BrokerDeploymentWrapper{}
		bdw.
			WithWait(true).
			WithBrokerClient(sw.BrokerClient).
			WithContext(ctx1).
			WithCustomImage(test.Config.BrokerImageName).
			WithName(DeployName)
	})

	ginkgo.BeforeEach(func() {}, 10)

	ginkgo.It("Deploy a broker instance and check that statistics endpoint works", func() {
		bdw.SetConsoleExposure(true)
		testStatistics()
	})

	ginkgo.It("Deploy a broker with console disabled and check that statistics endpoint works", func() {
		bdw.SetConsoleExposure(false)
		testStatistics()
	})

	ginkgo.It("Deploy a broker. By-default, statistics should be disabled.", func() {
		bdw.SetConsoleExposure(true)
		bdw.DeployBrokers(1)
		urls, err := bdw.GetExternalUrls(ExpectedUrl, 0)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		// there should be only single url in return in this case.
		url := fmt.Sprintf("%s://%s/%s/", Protocol, urls[0], AddressBit)
		resp, err := http.Get(url)
		gomega.Expect(err).To(gomega.HaveOccurred())
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusNotFound))
	})

})
