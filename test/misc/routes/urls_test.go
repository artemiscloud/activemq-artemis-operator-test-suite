package routes

import (
	"github.com/artemiscloud/activemq-artemis-operator-test-suite/pkg/bdw"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/framework"
)

var _ = ginkgo.Describe("RouteTests", func() {

	var (
		ctx1           *framework.ContextData
		brokerDeployer *bdw.BrokerDeploymentWrapper
	)

	const (
		// Should be available at all times
		ExpectedWsconsUrlPart = "wconsj"
		ExpectedAmqpUrlPart   = "amqp"
	)

	// PrepareNamespace after framework has been created
	ginkgo.JustBeforeEach(func() {
		ctx1 = sw.Framework.GetFirstContext()
		brokerDeployer = &bdw.BrokerDeploymentWrapper{}
		setEnv(ctx1, brokerDeployer)
	})
	//

	ginkgo.It("Deploy a broker instance to check default amqp url", func() {
		gomega.Expect(brokerDeployer.DeployBrokers(1)).To(gomega.BeNil(), "Broker delpoyment failed")
		_, err := brokerDeployer.GetExternalUrls(ExpectedAmqpUrlPart, 0)
		//URL should be created for this scenario
		gomega.Expect(err).NotTo(gomega.HaveOccurred(), "URL retrieval failed: %s", err)
	})

	ginkgo.It("Deploy a broker instance with wscons disabled", func() {
		brokerDeployer.WithConsoleExposure(false)
		gomega.Expect(brokerDeployer.DeployBrokers(1)).To(gomega.BeNil(), "Broker deployment failed")
		_, err := brokerDeployer.GetExternalUrls(ExpectedWsconsUrlPart, 0)
		//No URL should be created for this scenario
		gomega.Expect(err).To(gomega.HaveOccurred(), "Console URL has been created despite being disabled")
	})

	ginkgo.It("Deploy a broker instance with wscons enabled", func() {
		brokerDeployer.WithConsoleExposure(true)
		gomega.Expect(brokerDeployer.DeployBrokers(1)).To(gomega.BeNil(), "Broker deployment failed")
		_, err := brokerDeployer.GetExternalUrls(ExpectedWsconsUrlPart, 0)
		//URL should be created for this scenario
		gomega.Expect(err).NotTo(gomega.HaveOccurred(), "WSCons URL retrieval failed: %s", err)
	})

	ginkgo.It("Deploy a broker instance with wscons disabled, then enable it", func() {
		brokerDeployer.WithConsoleExposure(false)
		gomega.Expect(brokerDeployer.DeployBrokers(1)).To(gomega.BeNil(), "Broker deployment failed")
		_, err := brokerDeployer.GetExternalUrls(ExpectedWsconsUrlPart, 0)
		//No URL should be created for this scenario
		gomega.Expect(err).To(gomega.HaveOccurred(), "Console URL has been created despite being disabled")
		brokerDeployer.WithConsoleExposure(true)
		gomega.Expect(brokerDeployer.Update()).NotTo(gomega.HaveOccurred())
		_, err = brokerDeployer.GetExternalUrls(ExpectedWsconsUrlPart, 0)
		//URL should be created for this scenario
		gomega.Expect(err).NotTo(gomega.HaveOccurred(), "Console URL has not been created after configuration change")
	})

})
