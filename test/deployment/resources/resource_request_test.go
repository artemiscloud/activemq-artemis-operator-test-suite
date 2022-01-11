package resources

import (
	"github.com/artemiscloud/activemq-artemis-operator-test-suite/pkg/bdw"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/framework"
)

var _ = ginkgo.Describe("ResourceRequestsTests", func() {

	var (
		ctx1           *framework.ContextData
		brokerDeployer *bdw.BrokerDeploymentWrapper
	)

	ginkgo.JustBeforeEach(func() {
		ctx1 = sw.Framework.GetFirstContext()
		brokerDeployer = &bdw.BrokerDeploymentWrapper{}
		setEnv(ctx1, brokerDeployer)
	})

	ginkgo.It("CPU Request check", func() {
		expectedCPURequest := "500m" // half vCPU/Core/Hyperthread
		brokerDeployer.WithCPURequest(expectedCPURequest)
		deployBroker(brokerDeployer)
		pod := getPod(ctx1)
		actualCPURequest := pod.Spec.Containers[0].Resources.Requests.Cpu()
		gomega.Expect(expectedCPURequest).To(gomega.Equal(actualCPURequest.String()), "Expected CPU Request: %s, real: %s", expectedCPURequest, actualCPURequest.String())
	})

	ginkgo.It("Memory Request check", func() {
		expectedMemRequest := "512M"
		brokerDeployer.WithMemRequest(expectedMemRequest)
		deployBroker(brokerDeployer)
		pod := getPod(ctx1)
		actualMemRequest := pod.Spec.Containers[0].Resources.Requests.Memory()
		gomega.Expect(expectedMemRequest).To(gomega.Equal(actualMemRequest.String()), "Expected Memory limit: %s, real: %s", expectedMemLimit, actualMemLimit.String())
	})

	ginkgo.It("Memory Request check", func() {
		expectedMemRequest := "512M"
		brokerDeployer.WithMemRequest(expectedMemRequest)
		deployBroker(brokerDeployer)
		pod := getPod(ctx1)

		actualMemRequest := pod.Spec.Containers[0].Resources.Requests.Memory()
		gomega.Expect(expectedMemRequest).To(gomega.Equal(actualMemRequest.String()), "Expected Memory limit: %s, real: %s", expectedMemLimit, actualMemLimit.String())

		expectedMemRequest = "768M"
		brokerDeployer.WithMemRequest(expectedMemRequest)
		brokerDeployer.Update()
		actualMemRequest := pod.Spec.Containers[0].Resources.Requests.Memory()
		gomega.Expect(expectedMemRequest).To(gomega.Equal(actualMemRequest.String()), "Expected Update to Memory limit: %s, real: %s", expectedMemLimit, actualMemLimit.String())

	})
})
