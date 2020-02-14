package basic

import (
	"github.com/onsi/ginkgo"
)

var _ = ginkgo.Describe("DeploymentSingleBroker", func() {

	var (
		//ctx1 *framework.ContextData
		//brokerClient brokerclientset.Interface
	)


	// Initialize after framework has been created
	ginkgo.JustBeforeEach(func() {
	//	ctx1 = Framework.GetFirstContext()
	})
/*
	ginkgo.It("Deploy single broker instance", func() {
		//ctx1.OperatorMap[operators.OperatorTypeBroker].Namespace()
		err := DeployBrokers(ctx1, 1)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("Deploy double broker instances", func() {
		//ctx1.OperatorMap[operators.OperatorTypeBroker].Namespace()
		err := DeployBrokers(ctx1, 2)
		gomega.Expect(err).To(gomega.BeNil())
	})
*/
})
