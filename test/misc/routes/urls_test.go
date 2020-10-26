package routes

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/framework"
	bdw "gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/pkg/bdw"
	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/test"
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
		brokerDeployer.
			WithWait(true).
			WithBrokerClient(sw.BrokerClient).
			WithContext(ctx1).
			WithCustomImage(test.Config.BrokerImageName).
			WithName(DeployName).
			WithLts(!test.Config.NeedsLatestCR)
	})
	//

	ginkgo.It("Deploy a broker instance to check default amqp url", func() {
		gomega.Expect(brokerDeployer.DeployBrokers(1)).To(gomega.BeNil())
		_, err := brokerDeployer.GetExternalUrls(ExpectedAmqpUrlPart, 0)
		//URL should be created for this scenario
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	})

	ginkgo.It("Deploy a broker instance with wscons disabled", func() {
		brokerDeployer.WithConsoleExposure(false)
		gomega.Expect(brokerDeployer.DeployBrokers(1)).To(gomega.BeNil())
		_, err := brokerDeployer.GetExternalUrls(ExpectedWsconsUrlPart, 0)
		//No URL should be created for this scenario
		gomega.Expect(err).To(gomega.HaveOccurred())
	})

	ginkgo.It("Deploy a broker instance with wscons enabled", func() {
		brokerDeployer.WithConsoleExposure(true)
		gomega.Expect(brokerDeployer.DeployBrokers(1)).To(gomega.BeNil())
		_, err := brokerDeployer.GetExternalUrls(ExpectedWsconsUrlPart, 0)
		//URL should be created for this scenario
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	})

	ginkgo.It("Deploy a broker instance with wscons disabled, then enable it", func() {
		brokerDeployer.WithConsoleExposure(false)
		gomega.Expect(brokerDeployer.DeployBrokers(1)).To(gomega.BeNil())
		_, err := brokerDeployer.GetExternalUrls(ExpectedWsconsUrlPart, 0)
		//No URL should be created for this scenario
		gomega.Expect(err).To(gomega.HaveOccurred())
		brokerDeployer.WithConsoleExposure(true)
		gomega.Expect(brokerDeployer.Update()).NotTo(gomega.HaveOccurred())
		_, err = brokerDeployer.GetExternalUrls(ExpectedWsconsUrlPart, 0)
		//URL should be created for this scenario
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	})

})
