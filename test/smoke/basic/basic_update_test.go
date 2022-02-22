package basic

import (
	"strings"

	"github.com/artemiscloud/activemq-artemis-operator-test-suite/pkg/bdw"
	"github.com/artemiscloud/activemq-artemis-operator-test-suite/test"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/framework"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
		images := test.GetImages()
		imageName := decideImageName()
		imageArch := decideImageArch()
		CustomImage := ""
		for _, item := range images {
			if strings.HasPrefix(item.Name, imageName) && strings.HasSuffix(item.Name, imageArch) {
				CustomImage = item.Value
				break
			}
		}
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

func decideImageArch() string {
	name := ""
	if test.Config.PPC {
		name = "_ppc64le"
	} else if test.Config.IBMz {
		name = "_s390x"
	} else {
		//Problem: _s390x and _ppc here would still work.. Need an elegant solution
	}
	return name
}

func decideImageName() string {
	name := "RELATED_IMAGE_ActiveMQ_Artemis_Broker_Kubernetes"
	if test.Config.BrokerName != "amq-broker" {
		name = "RELATED_IMAGE_ActiveMQ_Artemis_Broker_Kubernetes"
	}
	return name
}

func getPod(ctx1 *framework.ContextData) *v1.Pod {
	kubeclient := ctx1.Clients.KubeClient
	podName := DeployName + PodNameSuffix
	pod, err := kubeclient.CoreV1().Pods(ctx1.Namespace).Get(podName, metav1.GetOptions{})
	gomega.Expect(err).To(gomega.BeNil())

	return pod
}
