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
		bdw  *test.BrokerDeploymentWrapper
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

	ginkgo.It("Deploy single broker, replace image with new one", func() {
		gomega.Expect(bdw.DeployBrokers(1)).To(gomega.BeNil())
		bdw.WithCustomImage(test.Config.BrokerImageOther)
		gomega.Expect(bdw.ChangeImage()).To(gomega.BeNil())
		gomega.Expect(bdw.VerifyImage(test.Config.BrokerImageOther)).NotTo(gomega.HaveOccurred())
	})

	ginkgo.It("Deploy single broker, scale down, replace image with new one, scale up", func() {
		gomega.Expect(bdw.DeployBrokers(1)).To(gomega.BeNil())
		gomega.Expect(bdw.Scale(0)).To(gomega.BeNil())
		bdw.WithCustomImage(test.Config.BrokerImageOther)
		gomega.Expect(bdw.ChangeImage()).To(gomega.BeNil())
		gomega.Expect(bdw.Scale(1)).To(gomega.BeNil())
		gomega.Expect(bdw.VerifyImage(test.Config.BrokerImageOther)).NotTo(gomega.HaveOccurred())
	})
})
