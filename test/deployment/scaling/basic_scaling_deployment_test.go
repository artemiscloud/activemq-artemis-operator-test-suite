package scaling

import (
	"github.com/artemiscloud/activemq-artemis-operator-test-suite/pkg/bdw"
	"github.com/artemiscloud/activemq-artemis-operator-test-suite/test"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/framework"
)

var _ = ginkgo.Describe("DeploymentScalingBroker", func() {

	var (
		ctx1           *framework.ContextData
		brokerDeployer *bdw.BrokerDeploymentWrapper
	)

	// PrepareNamespace after framework has been created
	ginkgo.JustBeforeEach(func() {
		ctx1 = sw.Framework.GetFirstContext()
		brokerDeployer = &bdw.BrokerDeploymentWrapper{}
		brokerDeployer.
			WithWait(true).
			WithBrokerClient(sw.BrokerClient).
			WithContext(ctx1).
			WithCustomImage(test.Config.BrokerImageName).
			WithName(DeployName).
			WithLts(!test.Config.NeedsLatestCR)

	})

	ginkgo.It("Deploy single broker instance and scale it to 4 replicas", func() {
		gomega.Expect(brokerDeployer.DeployBrokers(1)).To(gomega.BeNil(), "Broker deployment failed")
		gomega.Expect(brokerDeployer.Scale(4)).To(gomega.BeNil(), "Broker scaling to 4 failed")
	})

	ginkgo.It("Deploy multiple broker instances and scale it down to 1", func() {
		gomega.Expect(brokerDeployer.DeployBrokers(4)).To(gomega.BeNil(), "Broker deployment failed")
		gomega.Expect(brokerDeployer.Scale(1)).To(gomega.BeNil(), "Broker downscaling failed")
	})

	ginkgo.It("Deploy single broker instances and scale it down to 0", func() {
		gomega.Expect(brokerDeployer.DeployBrokers(1)).To(gomega.BeNil(), "Broker deployment failed")
		gomega.Expect(brokerDeployer.Scale(0)).To(gomega.BeNil(), "Broker downscaling to 0 failed")
	})

	ginkgo.It("Deploy zero broker instances and scale up to 1", func() {
		gomega.Expect(brokerDeployer.WithWait(false).DeployBrokers(0)).To(gomega.BeNil(), "Creation of an empty statefulset failed")
		gomega.Expect(brokerDeployer.WithWait(true).Scale(1)).To(gomega.BeNil(), "Scaling of empty stateful set to 1 failed")
	})

	ginkgo.It("Deploy single broker instance and scale up to max (16)", func() {
		gomega.Expect(brokerDeployer.WithWait(true).DeployBrokers(1)).To(gomega.BeNil(), "Broker deployment failed")
		gomega.Expect(brokerDeployer.WithWait(true).Scale(16)).To(gomega.BeNil(), "Broker scaling to 16 failed")
	})

})
