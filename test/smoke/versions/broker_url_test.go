package versions

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/framework"
	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/test"
	"time"
)

var _ = ginkgo.Describe("DeploymentScalingBroker", func() {

	var (
		ctx1 *framework.ContextData
		dw test.DeploymentWrapper
	)

	// PrepareNamespace after framework has been created
	ginkgo.JustBeforeEach(func() {
		ctx1 = Framework.GetFirstContext()
		dw = test.DeploymentWrapper{}.WithWait(true).WithBrokerClient(brokerClient).WithContext(ctx1).WithCustomImage(test.TestConfig.BrokerImageName)

	})
	// This tests various states in which containers could arrive. Malformed URL is different from valid URL and valid container URL.

	ginkgo.It("Define wrong (but valid) url for broker image, then replace with proper one", func() {
		gomega.Expect(dw.WithWait(false).WithCustomImage( "https://localhost/thing").DeployBrokers(1)).To(gomega.BeNil())
		time.Sleep(time.Duration(10)*time.Second)
		gomega.Expect(dw.WithWait(true).WithCustomImage(test.TestConfig.BrokerImageName).DeployBrokers(1)).To(gomega.BeNil())
		podLog, _ := ctx1.GetLogs("ex-aao-ss-0")
		gomega.Expect(podLog).To(gomega.ContainSubstring(test.TestConfig.BrokerVersion))
	})

	ginkgo.It("Define gibberish url for broker image then replace with proper one", func() {
		gomega.Expect(dw.WithCustomImage("gibberish://non-url-at-all").DeployBrokers( 1)).To(gomega.BeNil())
		time.Sleep(time.Duration(10)*time.Second)
		gomega.Expect(dw.WithCustomImage(test.TestConfig.BrokerImageName).ChangeImage())
		podLog, _ := ctx1.GetLogs("ex-aao-ss-0")
		gomega.Expect(podLog).To(gomega.ContainSubstring(test.TestConfig.BrokerVersion))
	})

	ginkgo.It("Define wrong image and replace it with broker", func() {
		gomega.Expect(dw.WithCustomImage(test.TestConfig.OperatorImageName).DeployBrokers( 1)).To(gomega.BeNil())
		time.Sleep(time.Duration(10)*time.Second)
		gomega.Expect(dw.WithCustomImage(test.TestConfig.BrokerImageName).ChangeImage())
		podLog, _ := ctx1.GetLogs("ex-aao-ss-0")
		gomega.Expect(podLog).To(gomega.ContainSubstring(test.TestConfig.BrokerVersion))
	})

	ginkgo.It("Define empty image url then replace it with broker", func() {
		gomega.Expect(dw.WithWait(false).WithCustomImage( "").DeployBrokers(1)).To(gomega.BeNil())
		time.Sleep(time.Duration(10)*time.Second)
		gomega.Expect(dw.WithWait(true).WithCustomImage(test.TestConfig.BrokerImageName).DeployBrokers(1)).To(gomega.BeNil())
		podLog, _ := ctx1.GetLogs("ex-aao-ss-0")
		gomega.Expect(podLog).To(gomega.ContainSubstring(test.TestConfig.BrokerVersion))
	})
})
