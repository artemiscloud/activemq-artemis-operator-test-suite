package basic

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/framework"
	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/test"
)

var _ = ginkgo.Describe("DeploymentSingleBroker", func() {

	var (
		ctx1 *framework.ContextData
		//brokerClient brokerclientset.Interface
	)


	// Initialize after framework has been created
	ginkgo.JustBeforeEach(func() {
		ctx1 = Framework.GetFirstContext()
	})

	ginkgo.It("Deploy single broker instance", func() {
		//ctx1.OperatorMap[operators.OperatorTypeBroker].Namespace()
		err := test.DeployBrokers(ctx1, 1, brokerClient)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("Deploy double broker instances", func() {
		//ctx1.OperatorMap[operators.OperatorTypeBroker].Namespace()
		err := test.DeployBrokers(ctx1, 2, brokerClient)
		gomega.Expect(err).To(gomega.BeNil())
	})

})
