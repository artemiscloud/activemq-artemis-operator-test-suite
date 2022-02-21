package bdw

import (
	"time"

	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/framework"
	"github.com/rh-messaging/shipshape/pkg/framework/log"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (bdw *BrokerDeploymentWrapper) WaitForBrokerSet(result int, originalSize int) {

	log.Logf("Timeout: %s", bdw.GetTimeout(max(result, originalSize)))

	err := framework.WaitForStatefulSet(bdw.ctx1.Clients.KubeClient,
		bdw.ctx1.Namespace,
		bdw.name+"-ss",
		result,
		time.Second*10, bdw.GetTimeout(max(result, originalSize)))
	gomega.Expect(err).To(gomega.BeNil(), "Deployment of broker failed: %s", err)
}

func (bdw *BrokerDeploymentWrapper) VerifyImage(target string) error {
	artemisCreated, err := bdw.brokerClient.BrokerV2alpha1().ActiveMQArtemises(bdw.ctx1.Namespace).Get(bdw.name, v1.GetOptions{})
	if err != nil {
		return err
	}
	gomega.Expect(artemisCreated.Spec.DeploymentPlan.Image).To(gomega.Equal(target))
	return nil
}

func (bdw *BrokerDeploymentWrapper) GetTimeout(count int) time.Duration {
	return time.Minute * time.Duration(5*bdw.timeoutMult*count)
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}
