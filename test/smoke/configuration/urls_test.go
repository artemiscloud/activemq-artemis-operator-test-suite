package configuration

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/framework"
	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/test"
)

var _ = ginkgo.Describe("RouteTests", func() {

	var (
		ctx1 *framework.ContextData
		bdw  *test.BrokerDeploymentWrapper
	)

	const (
		// Should be available at all times
		ExpectedWsconsUrlPart = "wconsj"
		ExpectedAmqpUrlPart   = "amqp"
	)

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
	//

	ginkgo.It("Deploy a broker instance to check default amqp url", func() {
		gomega.Expect(bdw.DeployBrokers(1)).To(gomega.BeNil())
		_, err := bdw.GetExternalUrls(ExpectedAmqpUrlPart, 0)
		//URL should be created for this scenario
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	})

	ginkgo.It("Deploy a broker instance with wscons disabled", func() {
		bdw.SetConsoleExposure(false)
		gomega.Expect(bdw.DeployBrokers(1)).To(gomega.BeNil())
		_, err := bdw.GetExternalUrls(ExpectedWsconsUrlPart, 0)
		//No URL should be created for this scenario
		gomega.Expect(err).To(gomega.HaveOccurred())
	})

	ginkgo.It("Deploy a broker instance with wscons enabled", func() {
		bdw.SetConsoleExposure(true)
		gomega.Expect(bdw.DeployBrokers(1)).To(gomega.BeNil())
		_, err := bdw.GetExternalUrls(ExpectedWsconsUrlPart, 0)
		//URL should be created for this scenario
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	})

	ginkgo.It("Deploy a broker instance with wscons disabled, then enable it", func() {
		bdw.SetConsoleExposure(false)
		gomega.Expect(bdw.DeployBrokers(1)).To(gomega.BeNil())
		_, err := bdw.GetExternalUrls(ExpectedWsconsUrlPart, 0)
		//No URL should be created for this scenario
		gomega.Expect(err).To(gomega.HaveOccurred())
		bdw.SetConsoleExposure(true)
		gomega.Expect(bdw.Update()).NotTo(gomega.HaveOccurred())
		_, err = bdw.GetExternalUrls(ExpectedWsconsUrlPart, 0)
		//URL should be created for this scenario
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	})

})
