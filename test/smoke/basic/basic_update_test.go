package basic

import (
	"context"
	"strings"

	"github.com/artemiscloud/activemq-artemis-operator-test-suite/pkg/bdw"
	"github.com/artemiscloud/activemq-artemis-operator-test-suite/test"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/framework"
	"github.com/rh-messaging/shipshape/pkg/framework/log"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	PPC  = "ppc64le"
	IBMZ = "s390x"
)

func determineInit() string {
	initName := "RELATED_IMAGE_ActiveMQ_Artemis_Broker_Init"
	imageArch := decideImageArch()
	CustomInit := ""
	images := test.GetImages()
	for _, item := range images {
		if strings.HasPrefix(item.Name, initName) && strings.HasSuffix(item.Name, imageArch) {
			if imageArch != "" {
				if strings.HasSuffix(item.Name, PPC) || strings.HasSuffix(item.Name, IBMZ) {
					CustomInit = item.Value
					break
				}
			} else {
				CustomInit = item.Value
				break
			}
		}
	}
	return CustomInit
}

func determineImage() string {
	images := test.GetImages()
	imageName := decideImageName()
	imageArch := decideImageArch()
	CustomImage := ""
	for _, item := range images {
		if strings.HasPrefix(item.Name, imageName) && strings.HasSuffix(item.Name, imageArch) {
			if imageArch != "" { // Also check lack of other architectures..
				if strings.HasSuffix(item.Name, PPC) || strings.HasSuffix(item.Name, IBMZ) {
					CustomImage = item.Value
					break
				}
			} else {
				CustomImage = item.Value
				break
			}
		}
	}
	return CustomImage
}

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
		CustomImage := determineImage()
		CustomInit := determineInit()
		log.Logf("Image: %s", CustomImage)
		brokerDeployer.WithCustomImage(CustomImage)
		brokerDeployer.WithCustomInit(CustomInit)
		err := brokerDeployer.DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil(), "Broker deployment failed")
		pod := getPod(ctx1)
		actualImage := pod.Spec.Containers[0].Image
		gomega.Expect(actualImage).To(gomega.Equal(CustomImage), "Image not updated after CR update")

	})
	// TODO: find init image as well
})

func decideImageArch() string {
	name := ""
	if test.Config.PPC {
		name = PPC
	} else if test.Config.IBMz {
		name = IBMZ
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
	pod, err := kubeclient.CoreV1().Pods(ctx1.Namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	gomega.Expect(err).To(gomega.BeNil())

	return pod
}
