package basic

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/framework"
	bdw2 "gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/pkg/bdw"
	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/test"
)

var _ = ginkgo.Describe("DeploymentScalingBroker", func() {

	var (
		ctx1 *framework.ContextData
		bdw  *bdw2.BrokerDeploymentWrapper
	)

	// PrepareNamespace after framework has been created
	ginkgo.JustBeforeEach(func() {
		ctx1 = sw.Framework.GetFirstContext()
		bdw = &bdw2.BrokerDeploymentWrapper{}
		bdw.
			WithWait(true).
			WithBrokerClient(sw.BrokerClient).
			WithContext(ctx1).
			WithCustomImage(test.Config.BrokerImageName).
			WithName(DeployName).
		    WithLts(!test.Config.NeedsV2)

	})

	ginkgo.It("Deploy single broker instance and scale it to 4 replicas", func() {
		gomega.Expect(bdw.DeployBrokers(1)).To(gomega.BeNil())
		gomega.Expect(bdw.Scale(4)).To(gomega.BeNil())
	})

	ginkgo.It("Deploy multiple broker instances and scale it down to 1", func() {
		gomega.Expect(bdw.DeployBrokers(4)).To(gomega.BeNil())
		gomega.Expect(bdw.Scale(1)).To(gomega.BeNil())
	})

	ginkgo.It("Deploy single broker instances and scale it down to 0", func() {
		gomega.Expect(bdw.DeployBrokers(1)).To(gomega.BeNil())
		gomega.Expect(bdw.Scale(0)).To(gomega.BeNil())
	})

	ginkgo.It("Deploy zero broker instances and scale up to 1", func() {
		gomega.Expect(bdw.WithWait(false).DeployBrokers(0)).To(gomega.BeNil())
		gomega.Expect(bdw.WithWait(true).Scale(1)).To(gomega.BeNil())
	})

	ginkgo.It("Deploy single broker instance and scale up to max (16)", func() {
		gomega.Expect(bdw.WithWait(false).DeployBrokers(1)).To(gomega.BeNil())
		gomega.Expect(bdw.WithWait(true).Scale(16)).To(gomega.BeNil())
	})

})
