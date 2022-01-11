package basic

import (
	"github.com/artemiscloud/activemq-artemis-operator-test-suite/pkg/bdw"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/framework"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	CustomImage = "registry.redhat.io/amq7/amq-broker-rhel8@sha256:c5f4c08e068b9721967cf7c7cbd9a9e93fb5e39b264dd13b653e99d8f3fa9e0e"
)

var _ = ginkgo.Describe("DeploymentUpdateTests", func() {

	var (
		ctx1 *framework.ContextData
		//brokerClient brokerclientset.Interface
		brokerDeployer *bdw.BrokerDeploymentWrapper
	)

	// PrepareNamespace after framework has been created
	ginkgo.JustBeforeEach(func() {
		ctx1 = sw.Framework.GetFirstContext()
		brokerDeployer = &bdw.BrokerDeploymentWrapper{}
		setEnv(ctx1, brokerDeployer)
	})

	ginkgo.It("CustomImageOverrideTest", func() {
		brokerDeployer.WithCustomImage(CustomImage)
		// TODO: extract this from operator.yaml
		err := brokerDeployer.DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil(), "Broker deployment failed")
		//TODO	// Also verify image from the ""broker"" instance
		pod := getPod(ctx1)
		actualImage := pod.Spec.Containers[0].Image
		gomega.Expect(actualImage).To(gomega.Equal(CustomImage), "Image not updated after CR update")

	})

})

func getPod(ctx1 *framework.ContextData) *v1.Pod {
	kubeclient := ctx1.Clients.KubeClient
	podName := DeployName + PodNameSuffix
	pod, err := kubeclient.CoreV1().Pods(ctx1.Namespace).Get(podName, metav1.GetOptions{})
	gomega.Expect(err).To(gomega.BeNil())

	return pod
}
