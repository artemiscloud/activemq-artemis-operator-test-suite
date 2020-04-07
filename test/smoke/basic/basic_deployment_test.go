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
		dw test.DeploymentWrapper
	)

	// PrepareNamespace after framework has been created
	ginkgo.JustBeforeEach(func() {
		ctx1 = Framework.GetFirstContext()
		dw = test.DeploymentWrapper{}.WithWait(true).
			WithBrokerClient(brokerClient).
			WithContext(ctx1).
			WithCustomImage(test.Config.BrokerImageName).
			WithName(DeployName)
	})

	ginkgo.It("Deploy single broker instance", func() {
		//ctx1.OperatorMap[operators.OperatorTypeBroker].Namespace()
		err := dw.DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("Deploy double broker instances", func() {
		//ctx1.OperatorMap[operators.OperatorTypeBroker].Namespace()
		err := dw.DeployBrokers(2)
		gomega.Expect(err).To(gomega.BeNil())
	})


})
