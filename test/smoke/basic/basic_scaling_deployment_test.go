package basic

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/framework"
)

var _ = ginkgo.Describe("DeploymentScalingBroker", func() {

	var (
		ctx1 *framework.ContextData
	)

	// Initialize after framework has been created
	ginkgo.JustBeforeEach(func() {
		ctx1 = Framework.GetFirstContext()
	})
/*
	ginkgo.It("Deploy single broker instance and scale it to 4 replicas", func() {
		gomega.Expect(DeployBrokers(ctx1, 1)).To(gomega.BeNil())
		gomega.Expect(Scale(ctx1,4)).To(gomega.BeNil())
	})

	ginkgo.It("Deploy multiple broker instances and scale it down to 1", func() {
		gomega.Expect(DeployBrokers(ctx1, 4)).To(gomega.BeNil())
		gomega.Expect(Scale(ctx1,1)).To(gomega.BeNil())
	})

	ginkgo.It("Deploy single broker instances and scale it down to 0", func() {
		gomega.Expect(DeployBrokers(ctx1, 1)).To(gomega.BeNil())
		gomega.Expect(Scale(ctx1,0)).To(gomega.BeNil())	})
*/
	ginkgo.It("Deploy zero broker instances and scale up to 1", func() {
		gomega.Expect(DeployBrokers(ctx1, 0)).To(gomega.BeNil())
		gomega.Expect(Scale(ctx1,1)).To(gomega.BeNil())	})

})
