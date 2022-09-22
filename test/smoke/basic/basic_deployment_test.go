package basic

import (
	"context"
	"fmt"
	"strings"

	"github.com/artemiscloud/activemq-artemis-operator-test-suite/pkg/bdw"
	"github.com/artemiscloud/activemq-artemis-operator-test-suite/test"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/framework"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	PodNameSuffix = "-ss-0"
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
		setEnv(ctx1, brokerDeployer)
	})

	ginkgo.It("Deploy single broker instance", func() {
		//ctx1.OperatorMap[operators.OperatorTypeBroker].Namespace()
		err := brokerDeployer.DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil(), "Broker deployment failed")
	})

	ginkgo.It("Deploy multiple broker ssets", func() {
		err := brokerDeployer.DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil(), "Broker deployment failed")
		err = brokerDeployer.WithName("secondset").DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil(), "Broker deployment failed")
		err = brokerDeployer.WithName("thirdset").DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil(), "Broker deployment failed")
		err = brokerDeployer.WithName("fourthset").DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil(), "Broker deployment failed")
	})

	ginkgo.It("Deploy multiple broker sets with different configuration", func() {
		err := brokerDeployer.DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil(), "Broker deployment failed")
		err = brokerDeployer.WithName("secondset").DeployBrokers(2)
		gomega.Expect(err).To(gomega.BeNil(), "Broker deployment failed")
		err = brokerDeployer.WithName("thirdset").WithMigration(true).WithPersistence(true).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil(), "Broker deployment failed")
		CustomImage := determineImage()
		CustomInit := determineInit()
		err = brokerDeployer.WithName("fourthset").WithMigration(false).WithPersistence(true).WithCustomImage(CustomImage).WithCustomInit(CustomInit).DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil(), "Broker deployment failed")
		cv1 := ctx1.Clients.KubeClient.CoreV1()
		pod, err := cv1.Pods(ctx1.Namespace).Get(context.TODO(), "fourthset-ss-0", v1.GetOptions{})
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		gomega.Expect(pod.Spec.Containers[0].Image).To(gomega.Equal(CustomImage), "Custom image not applied to pod")
		gomega.Expect(pod.Spec.InitContainers[0].Image).To(gomega.Equal(CustomInit), "Custom image not applied to pod")
	})
	ginkgo.It("Deploy double broker instances", func() {
		//ctx1.OperatorMap[operators.OperatorTypeBroker].Namespace()

		err := brokerDeployer.DeployBrokers(2)
		gomega.Expect(err).To(gomega.BeNil(), "Double Broker deployment failed")
	})

	ginkgo.It("Deploy with a bogus property", func() {
		brokerDeployer.AddProperty("abyrvalg", "bluepinspio")
		err := brokerDeployer.DeployBrokers(1)
		gomega.Expect(err).NotTo(gomega.HaveOccurred(), "broker not deployed")
		ss := brokerDeployer.GetStatefulSet()
		// Filename being in JAVA_OPTS ensures that its getting fed into container's cmdline
		initjavaopts := getEnvVarValue("JAVA_OPTS", ss.Spec.Template.Spec.InitContainers[0].Env)
		filename := strings.Split(initjavaopts, "=")[1]
		propertiesfile, err := brokerDeployer.GetFile("basic-ss-0", "basic-container", filename, sw.Framework.GetConfig())
		gomega.Expect(propertiesfile).To(gomega.ContainSubstring("abyrvalg=bluepinspio"), "properties file doesn't contain the expected string")
		//		propertiesfile, err := brokerDeployer.GetFile("basic-ss-0", "basic-container", "/home/jboss/amq-broker/etc/"
		// TODO: Properties broken? Currently it only verifies existing issue of wrong array type being supplied to CR
	})

	ginkgo.It("Deploy older broker version", func() {
		err := brokerDeployer.WithVersion("7.8.3").WithCustomImage("").DeployBrokers(1)
		gomega.Expect(err).NotTo(gomega.HaveOccurred(), "broker not deployed")
		ss := brokerDeployer.GetStatefulSet()
		images := test.GetImages()
		imageArch := decideImageArch()
		imagever := "783"
		imageName := ""
		initName := ""
		if imageArch == "" {
			imageName = fmt.Sprintf("%s_%s", decideImageName(), imagever)
			initName = fmt.Sprintf("RELATED_IMAGE_ActiveMQ_Artemis_Broker_Init%s", imagever)
		} else {
			imageName = fmt.Sprintf("%s_%s_%s", decideImageName(), imagever, imageArch)
			initName = fmt.Sprintf("RELATED_IMAGE_ActiveMQ_Artemis_Broker_Init_%s%s", imagever, imageArch)
		}
		gomega.Expect(ss.Spec.Template.Spec.Containers[0].Image).To(gomega.Equal(getEnvVarValue(imageName, images)), "wrong broker image used in actual SS")
		gomega.Expect(ss.Spec.Template.Spec.InitContainers[0].Image).To(gomega.Equal(getEnvVarValue(initName, images)), "wrong broker image used in actual SS")

	})
})

func getEnvVarValue(name string, env []corev1.EnvVar) string {
	for _, item := range env {
		if item.Name == name {
			return item.Value
		}
	}
	return ""
}
