package basic

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/framework"
	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/test"
)

var _ = ginkgo.Describe("DeploymentScalingBroker", func() {

	var (
		ctx1 *framework.ContextData
		dw test.DeploymentWrapper

	)

	// PrepareNamespace after framework has been created
	ginkgo.JustBeforeEach(func() {
		ctx1 = Framework.GetFirstContext()
		dw = test.DeploymentWrapper{}.WithWait(true).WithBrokerClient(brokerClient).WithContext(ctx1).WithCustomImage(test.TestConfig.BrokerImageName)
	})

	ginkgo.It("Deploy single broker instance and scale it to 4 replicas", func() {
		gomega.Expect(dw.DeployBrokers(1)).To(gomega.BeNil())
		gomega.Expect(dw.Scale(4)).To(gomega.BeNil())
	})

	ginkgo.It("Deploy multiple broker instances and scale it down to 1", func() {
		gomega.Expect(dw.DeployBrokers(4)).To(gomega.BeNil())
		gomega.Expect(dw.Scale(1)).To(gomega.BeNil())
	})

	ginkgo.It("Deploy single broker instances and scale it down to 0", func() {
		gomega.Expect(dw.DeployBrokers( 1)).To(gomega.BeNil())
		gomega.Expect(dw.Scale(0)).To(gomega.BeNil())	})

	ginkgo.It("Deploy zero broker instances and scale up to 1", func() {
		gomega.Expect(dw.WithWait(false).DeployBrokers(0)).To(gomega.BeNil())
		gomega.Expect(dw.WithWait(true).Scale(1)).To(gomega.BeNil())	})

	ginkgo.It("Deploy single broker instance and scale up to max (16)", func() {
		gomega.Expect(dw.WithWait(false).DeployBrokers(1)).To(gomega.BeNil())
		gomega.Expect(dw.WithWait(true).Scale(16)).To(gomega.BeNil())})

})
