package basic

import (
	"github.com/artemiscloud/activemq-artemis-operator-test-suite/pkg/bdw"
	"github.com/artemiscloud/activemq-artemis-operator-test-suite/test"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/framework"
)

var _ = ginkgo.Describe("DeploymentWithImageUpdates", func() {

	var (
		ctx1           *framework.ContextData
		brokerDeployer *bdw.BrokerDeploymentWrapper
	)

	// PrepareNamespace after framework has been created
	ginkgo.JustBeforeEach(func() {
		ctx1 = sw.Framework.GetFirstContext()
		brokerDeployer = &bdw.BrokerDeploymentWrapper{}
		setEnv(ctx1, brokerDeployer)
	})

	ginkgo.It("Deploy single broker, replace image with new one", func() {
		gomega.Expect(brokerDeployer.DeployBrokers(1)).To(gomega.BeNil())
		brokerDeployer.WithCustomImage(test.Config.BrokerImageOther)
		gomega.Expect(brokerDeployer.ChangeImage()).To(gomega.BeNil())
		gomega.Expect(brokerDeployer.VerifyImage(test.Config.BrokerImageOther)).NotTo(gomega.HaveOccurred())
	})

	ginkgo.It("Deploy single broker, scale down, replace image with new one, scale up", func() {
		gomega.Expect(brokerDeployer.DeployBrokers(1)).To(gomega.BeNil())
		gomega.Expect(brokerDeployer.Scale(0)).To(gomega.BeNil())
		brokerDeployer.WithCustomImage(test.Config.BrokerImageOther)
		gomega.Expect(brokerDeployer.ChangeImage()).To(gomega.BeNil())
		gomega.Expect(brokerDeployer.Scale(1)).To(gomega.BeNil())
		gomega.Expect(brokerDeployer.VerifyImage(test.Config.BrokerImageOther)).NotTo(gomega.HaveOccurred())
	})
})
