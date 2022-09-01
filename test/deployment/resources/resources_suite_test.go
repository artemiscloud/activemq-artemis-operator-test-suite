package resources

import (
	"context"
	"testing"

	"github.com/artemiscloud/activemq-artemis-operator-test-suite/pkg/bdw"
	"github.com/artemiscloud/activemq-artemis-operator-test-suite/test"
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/framework"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	PodNameSuffix = "-ss-0"
)

func TestResources(t *testing.T) {
	test.PrepareNamespace(t, "resources", "Resources Limitation Suite")
}

func TestMain(m *testing.M) {
	test.Initialize(m)
}

func setEnv(ctx1 *framework.ContextData, brokerDeployer *bdw.BrokerDeploymentWrapper) {
	brokerDeployer.WithWait(true).
		WithContext(ctx1).
		WithBrokerClient(sw.BrokerClient).
		WithCustomImage(test.Config.BrokerImageName).
		WithName(DeployName).
		WithLts(!test.Config.NeedsLatestCR).
		WithIncreasedTimeout(test.Config.TimeoutMultiplier)
}

func deployBroker(brokerDeployer *bdw.BrokerDeploymentWrapper) {
	err := brokerDeployer.DeployBrokers(1)
	gomega.Expect(err).To(gomega.BeNil())
}

func getPod(ctx1 *framework.ContextData) *v1.Pod {
	kubeclient := ctx1.Clients.KubeClient
	podName := DeployName + PodNameSuffix
	pod, err := kubeclient.CoreV1().Pods(ctx1.Namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	gomega.Expect(err).To(gomega.BeNil())

	return pod
}
