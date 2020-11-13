package basic

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/framework"
	bdw "github.com/artemiscloud/activemq-artemis-operator-test-suite/pkg/bdw"
	"github.com/artemiscloud/activemq-artemis-operator-test-suite/test"
)

var _ = ginkgo.Describe("DeploymentBasicTests", func() {

	var (
		ctx1 *framework.ContextData
		//brokerClient brokerclientset.Interface
		brokerDeployer *bdw.BrokerDeploymentWrapper
	)

	// PrepareNamespace after framework has been created
	ginkgo.JustBeforeEach(func() {
		ctx1 = sw.Framework.GetFirstContext()
		brokerDeployer = &bdw.BrokerDeploymentWrapper{}
		brokerDeployer.WithWait(true).
			WithBrokerClient(sw.BrokerClient).
			WithContext(ctx1).
			WithCustomImage(test.Config.BrokerImageName).
			WithName(DeployName).
			WithLts(!test.Config.NeedsLatestCR)
	})

	ginkgo.It("Deploy single broker instance", func() {
		//ctx1.OperatorMap[operators.OperatorTypeBroker].Namespace()
		err := brokerDeployer.DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("Deploy double broker instances", func() {
		//ctx1.OperatorMap[operators.OperatorTypeBroker].Namespace()
		err := brokerDeployer.DeployBrokers(2)
		gomega.Expect(err).To(gomega.BeNil())
	})

})
