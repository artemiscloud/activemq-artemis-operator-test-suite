package bdw

import (
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/framework"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

func (bdw *BrokerDeploymentWrapper) WaitForBrokerSet(result int, originalSize int) {
	err := framework.WaitForStatefulSet(bdw.ctx1.Clients.KubeClient,
		bdw.ctx1.Namespace,
		bdw.name+"-ss",
		result,
		time.Second*10, time.Minute*time.Duration(5*max(result, originalSize)))
	gomega.Expect(err).To(gomega.BeNil())
}

func (bdw *BrokerDeploymentWrapper) VerifyImage(target string) error {
	artemisCreated, err := bdw.brokerClient.BrokerV2alpha1().ActiveMQArtemises(bdw.ctx1.Namespace).Get(bdw.name, v1.GetOptions{})
	if err != nil {
		return err
	}
	gomega.Expect(artemisCreated.Spec.DeploymentPlan.Image).To(gomega.Equal(target))
	return nil
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}
