package basic

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/framework"
	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/test"
)

var _ = ginkgo.Describe("DeploymentWithImageUpdates", func() {

	var (
		ctx1 *framework.ContextData
		dw   *test.DeploymentWrapper
	)

	// PrepareNamespace after framework has been created
	ginkgo.JustBeforeEach(func() {
		ctx1 = Framework.GetFirstContext()
		dw = &test.DeploymentWrapper{}
		dw.
			WithWait(true).
			WithBrokerClient(brokerClient).
			WithContext(ctx1).
			WithCustomImage(test.Config.BrokerImageName).
			WithName(DeployName)
	})

	ginkgo.It("Deploy single broker, replace image with new one", func() {
		gomega.Expect(dw.DeployBrokers(1)).To(gomega.BeNil())
		dw.WithCustomImage(test.Config.BrokerImageOther)
		gomega.Expect(dw.ChangeImage()).To(gomega.BeNil())
		gomega.Expect(dw.VerifyImage(test.Config.BrokerImageOther)).NotTo(gomega.HaveOccurred())
	})

	ginkgo.It("Deploy single broker, scale down, replace image with new one, scale up", func() {
		gomega.Expect(dw.DeployBrokers(1)).To(gomega.BeNil())
		gomega.Expect(dw.Scale(0)).To(gomega.BeNil())
		dw.WithCustomImage(test.Config.BrokerImageOther)
		gomega.Expect(dw.ChangeImage()).To(gomega.BeNil())
		gomega.Expect(dw.Scale(1)).To(gomega.BeNil())
		gomega.Expect(dw.VerifyImage(test.Config.BrokerImageOther)).NotTo(gomega.HaveOccurred())
	})
})
