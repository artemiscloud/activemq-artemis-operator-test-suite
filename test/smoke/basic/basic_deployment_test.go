package basic

import (
	"context"

	"github.com/artemiscloud/activemq-artemis-operator-test-suite/pkg/bdw"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/framework"
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

})
