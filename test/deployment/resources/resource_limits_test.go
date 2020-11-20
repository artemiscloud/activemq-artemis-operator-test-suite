package resources

import (
	"github.com/artemiscloud/activemq-artemis-operator-test-suite/pkg/bdw"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/framework"
)

var _ = ginkgo.Describe("ResourceLimitsTests", func() {

	var (
		ctx1           *framework.ContextData
		brokerDeployer *bdw.BrokerDeploymentWrapper
	)

	ginkgo.JustBeforeEach(func() {
		ctx1 = sw.Framework.GetFirstContext()
		brokerDeployer = &bdw.BrokerDeploymentWrapper{}
		setEnv(ctx1, brokerDeployer)
	})

	ginkgo.It("CPU Limit check", func() {
		expectedCPULimit := "500m" // half vCPU/Core/Hyperthread
		brokerDeployer.WithCPULimit(expectedCPULimit)
		deployBroker(brokerDeployer)
		pod := getPod(ctx1)
		actualCPULimit := pod.Spec.Containers[0].Resources.Limits.Cpu()
		gomega.Expect(expectedCPULimit).To(gomega.Equal(actualCPULimit.String()))
	})

	ginkgo.It("Memory Limit check", func() {
		expectedMemLimit := "512M"
		brokerDeployer.WithMemLimit(expectedMemLimit)
		deployBroker(brokerDeployer)
		pod := getPod(ctx1)
		actualMemLimit := pod.Spec.Containers[0].Resources.Limits.Memory()
		gomega.Expect(expectedMemLimit).To(gomega.Equal(actualMemLimit.String()))
	})
})
